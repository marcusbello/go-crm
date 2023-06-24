package main

import (
	"flag"
	"fmt"
	"github.com/marcusbello/go-crm/internal/server"
	"github.com/marcusbello/go-crm/internal/server/storage/mem"
	"log"
)

// General service flags.
var (
	addr = flag.String("addr", "0.0.0.0:6742", "The address to run the service on.")
)

func main() {
	flag.Parse()

	// Setup for the service.
	store := mem.New()

	s, err := server.New(
		*addr,
		store,
		server.WithGRPCOpts(),
	)
	if err != nil {
		fmt.Printf("Unable to start grpc server, err: %v", err)
		panic(err)
	}

	done := make(chan error, 1)

	log.Println("Starting server at: ", *addr)
	go func() {
		defer close(done)
		done <- s.Start()
	}()

	log.Println("Server exited with error: ", <-done)
}
