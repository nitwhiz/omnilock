package server

import (
	"errors"
	"github.com/nitwhiz/omnilock/pkg/client"
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

		if timeout == 0 {
			result = s.LockTable.Lock(c, lockName)
		} else {
			result = s.LockTable.LockWithTimeout(c, lockName, timeout)
		}

		return result, nil
	}
}
