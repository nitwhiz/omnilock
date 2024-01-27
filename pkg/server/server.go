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
	clients       *ClientList
	cmdChan       chan *Command
	lockTable     *LockTable
	wg            *sync.WaitGroup
	cmdHandlers   map[string]CommandHandler
}

func New(ctx context.Context, opts ...Option) (*Server, error) {
	s := Server{
		clients:       NewClientList(),
		ctx:           ctx,
		keepAlive:     time.Second * 5,
		acceptorCount: runtime.NumCPU(),
		listenAddr:    "0.0.0.0:7194",
		cmdChan:       make(chan *Command),
		lockTable:     NewLockTable(ctx),
		wg:            &sync.WaitGroup{},
		cmdHandlers:   map[string]CommandHandler{},
	}

	for _, withOption := range opts {
		withOption(&s)
	}

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

	c := NewClient(conn, s.cmdChan)

	s.clients.Add(c)

	c.StartCommandLoop()

	s.lockTable.UnlockAllForClient(c)
}

func (s *Server) waitForShutdown() {
	<-time.After(time.Second)
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
