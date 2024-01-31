package server

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	keepAlive     time.Duration
	listener      net.Listener
	acceptorCount int
	listenAddr    string
	ctx           context.Context
	cmdChan       chan *Command
	lockTable     *LockTable
	clientCount   int
	wg            *sync.WaitGroup
	cmdHandlers   map[string]CommandHandler
}

func (s *Server) applyDefaults() {
	if s.keepAlive == 0 {
		WithKeepAlivePeriod(time.Second * 5)(s)
	}

	if s.acceptorCount == 0 {
		WithAcceptorCount(runtime.NumCPU())(s)
	}

	if s.listenAddr == "" {
		WithListenAddr("0.0.0.0:7194")(s)
	}
}

func New(ctx context.Context, opts ...Option) (*Server, error) {
	s := Server{
		ctx:           ctx,
		keepAlive:     0,
		acceptorCount: 0,
		listenAddr:    "",
		cmdChan:       make(chan *Command),
		lockTable:     NewLockTable(ctx),
		clientCount:   0,
		wg:            &sync.WaitGroup{},
		cmdHandlers:   map[string]CommandHandler{},
	}

	for _, withOption := range opts {
		withOption(&s)
	}

	s.applyDefaults()

	listenConfig := net.ListenConfig{
		KeepAlive: s.keepAlive,
	}

	listener, err := listenConfig.Listen(s.ctx, "tcp", s.listenAddr)

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
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	s.wg.Add(1)
	defer s.wg.Done()

	s.clientCount += 1
	defer func() {
		s.clientCount -= 1
	}()

	c := NewClient(conn, s.cmdChan)

	c.StartCommandLoop()

	s.lockTable.UnlockAllForClient(c)
}

func (s *Server) waitForShutdown() {
	<-s.ctx.Done()
	s.wg.Wait()
}

func (s *Server) Accept() {
	go s.startCommandListener()

	for i := 0; i < s.acceptorCount; i++ {
		go s.acceptor()
	}

	fmt.Println("Ready!")

	s.waitForShutdown()
}
