package server

import (
	"context"
	"github.com/nitwhiz/omnilock/v2/pkg/client"
	"github.com/nitwhiz/omnilock/v2/pkg/lock"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	options     *Options
	listener    net.Listener
	wg          *sync.WaitGroup
	ctx         context.Context
	locks       *lock.Table
	cmdHandlers map[string]CommandHandler
	cmdChan     chan *client.Command
}

func New(ctx context.Context, opts ...Option) (*Server, error) {
	s := Server{
		options:     &Options{},
		wg:          &sync.WaitGroup{},
		ctx:         ctx,
		locks:       lock.NewTable(),
		cmdHandlers: map[string]CommandHandler{},
		cmdChan:     make(chan *client.Command),
	}

	for _, withOption := range opts {
		withOption(s.options)
	}

	s.options.applyDefaults()

	listenConfig := net.ListenConfig{
		KeepAlive: s.options.keepAlive,
	}

	listener, err := listenConfig.Listen(ctx, "tcp", s.options.listenAddr)

	if err != nil {
		return nil, err
	}

	s.listener = listener

	s.initCommandHandlers()

	return &s, nil
}

func (s *Server) acceptor() {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			if s.ctx.Err() != nil {
				return
			}

			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	s.wg.Add(1)
	defer s.wg.Done()

	ctx, cancel := context.WithCancel(s.ctx)

	c := client.New(ctx, conn, s.cmdChan)

	clientId := c.GetID()

	defer log.Printf("Client #%d disonnected.\n", clientId)

	defer func() {
		c.Lock()
		defer c.Unlock()

		s.locks.UnlockAll(clientId)
	}()

	defer cancel()

	log.Printf("Client #%d connected from %s.\n", clientId, conn.RemoteAddr().String())

	c.ListenForCommands()
}

func (s *Server) wait() {
	<-s.ctx.Done()
	s.wg.Wait()
}

func (s *Server) Run() {
	go s.startCommandListener()

	for i := 0; i < s.options.acceptorCount; i++ {
		go s.acceptor()
	}

	log.Println("Ready to serve connections")

	<-time.After(time.Millisecond)

	s.wait()
}

func (s *Server) Write(c *client.Client, msg string) {
	_, _ = c.Write([]byte(msg+"\n"), s.options.clientTimeout)
}
