package server

import (
	"fmt"
	"strings"
)

type CommandError struct {
	Client  *Client
	Command string
	Argv    []string
	Message string
}

func (c *CommandError) Error() string {
	return fmt.Sprintf("[%d] `%s %s`: %s", c.Client.ID, c.Command, strings.Join(c.Argv, " "), c.Message)
}
