package client

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"rpc/pkg/codec"
	"rpc/pkg/protocol"
	"sync"
)

type Call struct {
	ID            uint64
	ServiceMethod string
	Args          any
	Reply         any
	Error         error
	Done          chan *Call
}

type Client struct {
	conn      net.Conn
	codec     codec.Codec
	seq       uint64
	pending   map[uint64]*Call
	closing   bool
	mutex     sync.Mutex
	sendMutex sync.Mutex
}

func Dial(address string, codec codec.Codec) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	client := &Client{
		conn:    conn,
		codec:   codec,
		pending: make(map[uint64]*Call),
		seq:     1,
	}
	go client.receive()
	return client, nil
}

func (c *Client) send(call *Call) error {
	c.mutex.Lock()

	if c.closing {
		return errors.New("connection is closing")
	}

	call.ID = c.seq
	c.seq++
	c.pending[call.ID] = call

	c.mutex.Unlock()

	serializedArgs, err := c.codec.Encode(call.Args)
	if err != nil {
		c.removeCall(call.ID)
		return err
	}

	// Serialize args
	request := &protocol.Request{
		ID:            call.ID,
		ServiceMethod: call.ServiceMethod,
		Args:          serializedArgs,
	}

	requestBytes, err := c.codec.Encode(request)
	if err != nil {
		c.removeCall(call.ID)
		return err
	}

	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()

	length := uint32(len(requestBytes))
	if err := binary.Write(c.conn, binary.BigEndian, length); err != nil {
		c.removeCall(call.ID)
		return err
	}
	_, err = c.conn.Write(requestBytes)
	if err != nil {
		c.removeCall(call.ID)
		return err
	}
	return nil
}

func (c *Client) removeCall(id uint64) {
	c.mutex.Lock()
	delete(c.pending, id)
	c.mutex.Unlock()
}

func (c *Client) receive() {
	for {
		// 1. Read from connection
		// 2. Deserialize into Response
		// 3. Find the call
		// 4. Deliever the response
		// 5. Remove from pending map

		var length uint32
		if err := binary.Read(c.conn, binary.BigEndian, &length); err != nil {
			// Connection error, terminate all pending calls
			c.terminatePendingCalls(err)
			return
		}

		data := make([]byte, length)
		if _, err := io.ReadFull(c.conn, data); err != nil {
			c.terminatePendingCalls(err)
			return
		}

		var resp protocol.Response
		if err := c.codec.Decode(data, &resp); err != nil {
			c.terminatePendingCalls(err)
			return
		}

		c.mutex.Lock()
		call := c.pending[resp.ID]
		delete(c.pending, resp.ID)
		c.mutex.Unlock()

		if call == nil {
			continue
		}

		if resp.Err != "" {
			call.Error = errors.New(resp.Err)
		} else {
			if len(resp.Data) > 0 {
				err := c.codec.Decode(resp.Data, call.Reply)
				if err != nil {
					call.Error = err
				}
			}
		}

		// send the call on its Done channel, wake up anyone waiting on that channel
		call.Done <- call
	}
}

func (c *Client) terminatePendingCalls(err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, call := range c.pending {
		call.Error = err
		call.Done <- call
	}
	c.pending = make(map[uint64]*Call)
}

func (c *Client) Call(serviceMethod string, args any, reply any) error {
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          make(chan *Call, 1),
	}

	if err := c.send(call); err != nil {
		return err
	}

	receivedCall := <-call.Done
	return receivedCall.Error
}

func (c *Client) Go(serviceMethod string, args any, reply any, done chan *Call) *Call {
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}

	if call.Done == nil {
		call.Done = make(chan *Call, 1)
	}

	if err := c.send(call); err != nil {
		call.Error = err
		call.Done <- call
	}

	return call
}

func (c *Client) Close() error {
	c.mutex.Lock()
	if c.closing {
		c.mutex.Unlock()
		return errors.New("connection already closing")
	}
	c.closing = true
	c.mutex.Unlock()
	return c.conn.Close()
}
