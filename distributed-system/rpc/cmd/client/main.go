package main

import (
	"log"
	"rpc/pkg/client"
	"rpc/pkg/codec"
)

type Args struct {
	A, B int
}

type Reply struct {
	Result int
}

func main() {
	c, err := client.Dial("localhost:9000", codec.NewGobCodec())
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	args := &Args{A: 10, B: 20}
	reply := &Reply{}

	err = c.Call("Math.Add", args, reply)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Result: %d + %d = %d", args.A, args.B, reply.Result)
}
