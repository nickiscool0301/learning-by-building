package main

import (
	"fmt"
	"map-reduce/internal/worker"
	"os"
)

func main() {
	id := os.Args[1]
	w, err := worker.New(id, "localhost:9000")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Worker %s starting...\n", id)
	w.Run()
}
