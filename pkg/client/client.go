package client

import (
	"bufio"
	"context"
	"github.com/nitwhiz/omnilock/pkg/id"
	"net"
	"sync"
	"time"
)

type Client struct {
	id      uint64
	ctx     context.Context
	conn    net.Conn
	reader  *bufio.Reader
	cmdChan chan<- *Command
	mu      *sync.Mutex
}

func New(ctx context.Context, conn net.Conn, cmdChan chan<- *Command) *Client {
	c := Client{
		id:      id.Next(),
		ctx:     ctx,
		conn:    conn,
		reader:  bufio.NewReader(conn),
		cmdChan: cmdChan,
		mu:      &sync.Mutex{},
	}

	return &c
}

func (c *Client) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Client) Lock() {
	c.mu.Lock()
}

func (c *Client) Unlock() {
	c.mu.Unlock()
}

func (c *Client) GetID() uint64 {
	return c.id
}

func (c *Client) GetContext() context.Context {
	return c.ctx
}

func (c *Client) Write(b []byte, timeout time.Duration) (int, error) {
	err := c.conn.SetWriteDeadline(time.Now().Add(timeout))

	if err != nil {
		return 0, err
	}

	return c.conn.Write(b)
}

func (c *Client) waitForCommand() bool {
	cmdString, err := c.reader.ReadString('\n')

	if err != nil {
		return false
	}

	select {
	case <-c.ctx.Done():
		return false
	case c.cmdChan <- &Command{
		Client: c,
		Cmd:    cmdString[:len(cmdString)-1],
	}:
		break
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
