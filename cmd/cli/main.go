package main

import (
	"context"
	"github.com/nitwhiz/omnilock/pkg/server"
	"log"
)

func main() {
	s, err := server.New(context.Background())

	if err != nil {
		log.Println(err)
		return
	}

	s.Run()
}
