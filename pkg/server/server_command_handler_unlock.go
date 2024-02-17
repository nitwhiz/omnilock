package server

import (
	"errors"
	"github.com/nitwhiz/omnilock/pkg/client"
	"log"
)

func UnlockHandler(s *Server) CommandHandler {
	return func(c *client.Client, argv ...string) (result bool, err error) {
		if len(argv) == 0 {
			return false, errors.New("not enough arguments")
		}

		lockName := argv[0]

		result = s.locks.Unlock(lockName, c.GetID())

		log.Printf("Client #%d requested unlock '%s': %v.\n", c.GetID(), lockName, result)

		return result, nil
	}
}
