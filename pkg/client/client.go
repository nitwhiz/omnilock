package client

import (
	"bufio"
	"context"
	"github.com/nitwhiz/omnilock/pkg/id"
	"net"
)

type Client struct {
	id      uint64
	ctx     context.Context
	conn    net.Conn
	reader  *bufio.Reader
	cmdChan chan<- *Command
}

func New(ctx context.Context, conn net.Conn, cmdChan chan<- *Command) *Client {
	c := Client{
		id:      id.Next(),
		ctx:     ctx,
		conn:    conn,
		reader:  bufio.NewReader(conn),
		cmdChan: cmdChan,
	}

	return &c
}

func (c *Client) GetID() uint64 {
	return c.id
}

func (c *Client) GetContext() context.Context {
	return c.ctx
}

func (c *Client) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Client) waitForCommand() bool {
	cmdString, err := c.reader.ReadString('\n')

	if err != nil {
		return false
	}

	c.cmdChan <- &Command{
		Client: c,
		Cmd:    cmdString[:len(cmdString)-1],
	}

	return true
}

// ListenForCommands blocks until the client disconnects
func (c *Client) ListenForCommands() {
	defer func() {
		_ = c.conn.Close()
	}()

	for {
		if !c.waitForCommand() {
			return
		}
	}
}
