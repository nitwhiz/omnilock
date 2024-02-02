package server

import (
	"errors"
	"github.com/nitwhiz/omnilock/pkg/client"
)

func TryLockHandler(s *Server) CommandHandler {
	return func(c *client.Client, argv ...string) (result bool, err error) {
		if len(argv) == 0 {
			return false, errors.New("not enough arguments")
		}

		return s.LockTable.TryLock(c, argv[0]), nil
	}
}
