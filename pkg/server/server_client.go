package server

import (
	"bufio"
	"net"
)

type Client struct {
	ID      uint64
	conn    net.Conn
	reader  *bufio.Reader
	cmdChan chan<- *Command
}

func NewClient(conn net.Conn, cmdChan chan<- *Command) *Client {
	return &Client{
		ID:      0,
		conn:    conn,
		reader:  bufio.NewReader(conn),
		cmdChan: cmdChan,
	}
}

func (c *Client) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Client) waitForCommand() bool {
	cmdString, err := c.reader.ReadString('\n')

	if err != nil {
		return false
	}

	c.cmdChan <- c.NewCommand(cmdString[:len(cmdString)-1])

	return true
}

// StartCommandLoop blocks until the client disconnects
func (c *Client) StartCommandLoop() {
	defer func() {
		_ = c.conn.Close()
	}()

	for {
		if !c.waitForCommand() {
			return
		}
	}
}
