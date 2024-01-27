package server

const ResponseSuccess = "success"
const ResponseFailed = "failed"

func (s *Server) writeToClient(c *Client, r string) {
	_, _ = c.Write([]byte(r + "\n"))
}

func (s *Server) writeResponse(c *Client, success bool) {
	if success {
		s.writeToClient(c, ResponseSuccess)
	} else {
		s.writeToClient(c, ResponseFailed)
	}
}
