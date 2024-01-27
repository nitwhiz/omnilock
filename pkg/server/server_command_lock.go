package server

import (
	"strconv"
	"time"
)

func (s *Server) handleLockCommand(c *Client, argv ...string) error {
	argc := len(argv)

	if argc == 0 {
		s.writeResponse(c, false)

		return &CommandError{
			Client:  c,
			Command: CmdLock,
			Argv:    argv,
			Message: "not enough arguments",
		}
	}

	lockName := argv[0]
	var timeout time.Duration

	if argc >= 2 {
		timeoutStr := argv[1]

		timeoutInt, err := strconv.Atoi(timeoutStr)

		if err != nil {
			return &CommandError{
				Client:  c,
				Command: CmdLock,
				Argv:    argv,
				Message: "unable to parse timeout",
			}
		}

		timeout = time.Millisecond * time.Duration(timeoutInt)
	}

	result := false

	if timeout == 0 {
		result = s.lockTable.Lock(c, lockName)
	} else {
		result = s.lockTable.LockWithTimeout(c, lockName, timeout)
	}

	s.writeResponse(c, result)

	return nil
}
