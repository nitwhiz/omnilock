package server

import (
	"fmt"
	"strings"
)

type CommandHandler func(c *Client, argv ...string) error

func (s *Server) addCommandHandler(cmd string, handler CommandHandler) {
	s.cmdHandlers[cmd] = handler
}

func (s *Server) initCommandHandlers() {
	s.addCommandHandler(CmdLock, s.handleLockCommand)
	s.addCommandHandler(CmdTryLock, s.handleTryLockCommand)
	s.addCommandHandler(CmdUnlock, s.handleUnlockCommand)
}

func (s *Server) handleCommand(cmd *Command) {
	argv := strings.Split(cmd.Command, " ")

	cmdHandler, ok := s.cmdHandlers[argv[0]]

	if !ok {
		fmt.Println("Unknown Command")
		return
	}

	if err := cmdHandler(cmd.Client, argv[1:]...); err != nil {
		fmt.Println(err)
	}
}
