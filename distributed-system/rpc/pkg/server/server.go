package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"rpc/pkg/codec"
	"rpc/pkg/protocol"
	"strings"
	"sync"
)

type Server struct {
	addr     string
	codec    codec.Codec
	services map[string]any
	mu       sync.RWMutex
}

func NewServer(addr string, codec codec.Codec) *Server {
	return &Server{
		addr:     addr,
		codec:    codec,
		services: make(map[string]any),
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

	serviceName, methodName, err := s.parseServiceMethod(req.ServiceMethod)
	if err != nil {
		s.writeResponse(conn, &protocol.Response{
			ID:   req.ID,
			Data: nil,
			Err:  err.Error(),
		})
		return
	}

	replyData, err := s.CallMethod(serviceName, methodName, req.Args)

	resp := &protocol.Response{
		ID:   req.ID,
		Data: replyData,
		Err:  "",
	}

	if err != nil {
		resp.Err = err.Error()
	}

	if err := s.writeResponse(conn, resp); err != nil {
		log.Printf("write response error: %v", err)
		return
	}

	log.Printf("response sent: ID=%d", req.ID)
}

func (s *Server) Register(name string, service any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.services[name]; exists {
		return fmt.Errorf("service %s already registered", name)
	}

	s.services[name] = service
	log.Printf("service %s registered", name)
	return nil
}

func (s *Server) parseServiceMethod(serviceMethod string) (string, string, error) {
	parts := strings.Split(serviceMethod, ".")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid service method: %s", serviceMethod)
	}
	return parts[0], parts[1], nil
}

func (s *Server) CallMethod(serviceName, methodName string, args []byte) ([]byte, error) {
	s.mu.RLock()
	service := s.services[serviceName]
	s.mu.RUnlock()

	if service == nil {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	// Get Method using reflection
	serviceValue := reflect.ValueOf(service)
	method := serviceValue.MethodByName(methodName)

	if !method.IsValid() {
		return nil, fmt.Errorf("method %s not found in service %s", methodName, serviceName)
	}

	// get method type
	methodType := method.Type()

	if methodType.NumIn() != 2 {
		return nil, fmt.Errorf("method %s has invalid number of input parameters", methodName)
	}

	argsType := methodType.In(0)  // First parameter type
	replyType := methodType.In(1) // Second parameter type

	argsValue := reflect.New(argsType.Elem())   // Create *Args
	replyValue := reflect.New(replyType.Elem()) // Create *Reply

	if err := s.codec.Decode(args, argsValue.Interface()); err != nil {
		return nil, fmt.Errorf("decode args error: %w", err)
	}

	returnValues := method.Call([]reflect.Value{argsValue, replyValue})
	var err error
	if len(returnValues) > 0 && !returnValues[0].IsNil() {
		err = returnValues[0].Interface().(error)
	}
	if err != nil {
		return nil, err
	}

	return s.codec.Encode(replyValue.Interface())
}
