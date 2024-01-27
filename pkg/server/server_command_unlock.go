package server

func (s *Server) handleUnlockCommand(c *Client, argv ...string) error {
	if len(argv) == 0 {
		s.writeResponse(c, false)

		return &CommandError{
			Client:  c,
			Command: CmdUnlock,
			Argv:    argv,
			Message: "not enough arguments",
		}
	}

	lockName := argv[0]

	result := s.lockTable.Unlock(c, lockName)

	s.writeResponse(c, result)

	return nil
}
