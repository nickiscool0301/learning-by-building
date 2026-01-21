package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"rpc/pkg/codec"
	"rpc/pkg/protocol"
)

type Server struct {
	addr  string
	codec codec.Codec
}

func NewServer(addr string, codec codec.Codec) *Server {
	return &Server{
		addr:  addr,
		codec: codec,
	}
}

func (s *Server) Start() error {
	// Implementation for starting the server
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) readRequest(conn net.Conn) (*protocol.Request, error) {
	var length uint32

	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("read length error: %w", err)
	}

	data := make([]byte, length)

	if _, err := io.ReadFull(conn, data); err != nil {
		log.Printf("read data error: %v", err)
		return nil, fmt.Errorf("read data error: %w", err)
	}

	var req protocol.Request
	if err := s.codec.Decode(data, &req); err != nil {
		return nil, fmt.Errorf("decode request error: %w", err)
	}

	return &req, nil
}

func (s *Server) writeResponse(conn net.Conn, resp *protocol.Response) error {
	respData, err := s.codec.Encode(resp)

	if err != nil {
		return fmt.Errorf("encode response error: %w", err)
	}

	// Write length
	if err := binary.Write(conn, binary.BigEndian, uint32(len(respData))); err != nil {
		return fmt.Errorf("write length error: %w", err)
	}

	if _, err := conn.Write(respData); err != nil {
		return fmt.Errorf("write response error: %w", err)
	}
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	/*
		- read bytes from connection: need to know how many bytes to read
		- decode into a Request
		- Log for debug
		- create Reponse
		- encode and write reponse
		- close connection
	*/
	defer conn.Close()

	req, err := s.readRequest(conn)
	if err != nil {
		log.Printf("read request error: %v", err)
		return
	}

	log.Printf("received request: %+v", req)

	resp := &protocol.Response{
		ID:   req.ID,
		Data: nil,
		Err:  "",
	}

	if err := s.writeResponse(conn, resp); err != nil {
		log.Printf("write response error: %v", err)
		return
	}
}
