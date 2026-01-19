package rpc

import (
	"encoding/gob"
	"errors"
	"net"
	"time"
)

type Client struct {
	address    string
	timeout    time.Duration
	maxRetries int
}

// constructor
func NewClient(address string) *Client {
	return &Client{
		address:    address,
		timeout:    5 * time.Second,
		maxRetries: 3,
	}
}

func (c *Client) Send(msg Message) (*Message, error) {
	var lastErr error

	for attemp := 0; attemp < c.maxRetries; attemp++ {
		resp, err := c.sendOnce(msg)
		if err == nil {
			return resp, nil
		}
		lastErr = err

		// Exponential backoff
		time.Sleep(time.Duration(1<<attemp) * 100 * time.Millisecond)
	}
	return nil, errors.New("max retries exceeded: " + lastErr.Error())
}

func (c *Client) sendOnce(msg Message) (*Message, error) {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(c.timeout))

	// Send
	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(&msg); err != nil {
		return nil, err
	}

	// Receive
	decoder := gob.NewDecoder(conn)
	var resp Message
	if err := decoder.Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
