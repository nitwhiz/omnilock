package main

import (
	"context"
	"fmt"
	"github.com/nitwhiz/omnilock/pkg/prom"
	"github.com/nitwhiz/omnilock/pkg/server"
)

func main() {
	s, err := server.New(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Starting metrics server")

	go func() {
		if err := prom.Listen(s); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}()

	fmt.Println("Starting lock server")

	s.Accept()
}
