package server

const CmdLock = "lock"
const CmdTryLock = "trylock"
const CmdUnlock = "unlock"

type Command struct {
	Command string
	Client  *Client
}

func (c *Client) NewCommand(cmd string) *Command {
	return &Command{
		Command: cmd,
		Client:  c,
	}
}

func (s *Server) startCommandListener() {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		select {
		case c := <-s.cmdChan:
			s.handleCommand(c)
			break
		case <-s.ctx.Done():
			_ = s.listener.Close()

			return
		}
	}
}
