package server

func (s *Server) handleTryLockCommand(c *Client, argv ...string) error {
	if len(argv) == 0 {
		s.writeResponse(c, false)

		return &CommandError{
			Client:  c,
			Command: CmdTryLock,
			Argv:    argv,
			Message: "not enough arguments",
		}
	}

	lockName := argv[0]

	result := s.lockTable.TryLock(c, lockName)

	s.writeResponse(c, result)

	return nil
}
