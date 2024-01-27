package main

import (
	"context"
	"fmt"
	"github.com/nitwhiz/omnilock/pkg/server"
)

func main() {
	s, err := server.New(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Starting server")

	s.Accept()
}
