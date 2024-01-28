package server

func (s *Server) GetCurrentLockCount() int {
	return s.lockTable.Count()
}
