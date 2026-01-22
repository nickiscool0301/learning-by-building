package main

import (
	"log"
	"rpc/pkg/codec"
	"rpc/pkg/server"
)

type Math struct{}

type Args struct {
	A, B int
}

type Reply struct {
	Result int
}

func (m *Math) Add(args *Args, reply *Reply) error {
	reply.Result = args.A + args.B
	log.Printf("Add called: %d + %d = %d", args.A, args.B, reply.Result)
	return nil
}

func main() {

	// Create server

	server := server.NewServer(":9000", codec.NewGobCodec())

	// Register the Math service
	if err := server.Register("Math", &Math{}); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting RPC server on :9000")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
