package server

func (s *Server) GetCurrentLockCount() int {
	return s.lockTable.Count()
}

func (s *Server) GetCurrentClientCount() int {
	return s.clientCount
}
