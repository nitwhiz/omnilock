package server

import (
	"context"
	"errors"
	"github.com/nitwhiz/omnilock/v2/pkg/client"
	"log"
	"strconv"
	"time"
)

func LockHandler(s *Server) CommandHandler {
	return func(c *client.Client, argv ...string) (result bool, err error) {
		argc := len(argv)

		if argc == 0 {
			return false, errors.New("not enough arguments")
		}

		lockName := argv[0]
		var timeout time.Duration

		if argc >= 2 {
			timeoutStr := argv[1]

			timeoutInt, err := strconv.Atoi(timeoutStr)

			if err != nil {
				return false, errors.New("unable to parse timeout")
			}

			timeout = time.Millisecond * time.Duration(timeoutInt)
		}

		result = false

		if timeout == 0 {
			result = s.locks.Lock(s.ctx, lockName, c.GetID())
		} else {
			ctx, cancel := context.WithTimeout(s.ctx, timeout)
			defer cancel()

			result = s.locks.Lock(ctx, lockName, c.GetID())
		}

		log.Printf("Client #%d requested lock '%s': %v.\n", c.GetID(), lockName, result)

		return result, nil
	}
}
