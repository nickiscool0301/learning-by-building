package rpc

import (
	"encoding/gob"
	"net"
	"sync"
	"time"
)

type Server struct {
	listener net.Listener
	handlers map[MessageType]HandlerFunc
	mu       sync.RWMutex
}

type HandlerFunc func(payload []byte) (MessageType, []byte, error)

func NewServer(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: listener,
		handlers: make(map[MessageType]HandlerFunc),
	}, nil
}

func (s *Server) Register(msgType MessageType, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[msgType] = handler
}

func (s *Server) Serve() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// Decode incoming message
	decoder := gob.NewDecoder(conn)
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		return
	}

	// Find handler
	s.mu.RLock()
	handler, ok := s.handlers[msg.Type]
	s.mu.RUnlock()

	// Build response
	var resp Message
	if !ok {
		resp = Message{Type: MsgAck, Payload: nil}
	} else {
		respType, respPayload, err := handler(msg.Payload)
		if err != nil {
			resp = Message{Type: MsgAck, Payload: nil}
		} else {
			resp = Message{Type: respType, Payload: respPayload}
		}
	}

	// Send response back
	encoder := gob.NewEncoder(conn)
	encoder.Encode(&resp)
}
