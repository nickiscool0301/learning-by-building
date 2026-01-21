package main

import (
	"fmt"
	"map-reduce/internal/coordinator"
	"map-reduce/internal/rpc"
	"os"
	"path/filepath"
)

func main() {
	// init coordinator
	c := coordinator.New()

	// create RPC server
	server, err := rpc.NewServer(":9000")
	if err != nil {
		panic(err)
	}

	// Register handlers
	c.RegisterHandlers(server)

	// Add test
	files, _ := filepath.Glob("test-data/*txt")
	for _, file := range files {
		content, _ := os.ReadFile(file)
		taskID := c.AddTask(file, string(content))
		fmt.Printf("Added task %s for file %s\n", taskID, file)
	}

	fmt.Printf("Added %d tasks\n", len(files))
	fmt.Println("Coordinator is running on :9000")

	server.Serve()

}
