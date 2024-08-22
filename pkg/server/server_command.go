package server

import (
	"github.com/nitwhiz/omnilock/v2/pkg/client"
	"strings"
)

func (s *Server) addCommandHandler(cmdName string, handler CommandHandler) {
	s.cmdHandlers[cmdName] = handler
}

func (s *Server) initCommandHandlers() {
	s.addCommandHandler("lock", LockHandler(s))
	s.addCommandHandler("trylock", TryLockHandler(s))
	s.addCommandHandler("unlock", UnlockHandler(s))
}

func (s *Server) handleCommand(cmd *client.Command) {
	s.wg.Add(1)
	defer s.wg.Done()

	cmd.Client.Lock()
	defer cmd.Client.Unlock()

	select {
	case <-cmd.Client.Done():
		return
	default:
		break
	}

	argv := strings.Split(cmd.Cmd, " ")

	if len(argv) < 1 {
		s.Write(cmd.Client, "error: missing command")
		return
	}

	cmdHandler, ok := s.cmdHandlers[argv[0]]

	if !ok {
		s.Write(cmd.Client, "error: unknown command")
		return
	}

	result, err := cmdHandler(cmd.Client, argv[1:]...)

	if err != nil {
		s.Write(cmd.Client, "error: "+err.Error())
		return
	}

	if result {
		s.Write(cmd.Client, "success")
	} else {
		s.Write(cmd.Client, "failed")
	}
}

func (s *Server) startCommandListener() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			_ = s.listener.Close()
			return
		case c := <-s.cmdChan:
			go s.handleCommand(c)
			break
		}
	}
}
