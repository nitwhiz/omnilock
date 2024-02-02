package server

import (
	"context"
	"fmt"
	"github.com/nitwhiz/omnilock/pkg/client"
	"github.com/nitwhiz/omnilock/pkg/locking"
	"net"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	wg            *sync.WaitGroup
	keepAlive     time.Duration
	listenAddr    string
	listener      net.Listener
	acceptorCount int
	ctx           context.Context
	LockTable     *locking.LockTable
	cmdChan       chan *client.Command
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
		wg:            &sync.WaitGroup{},
		keepAlive:     0,
		listenAddr:    "",
		listener:      nil,
		acceptorCount: 0,
		ctx:           ctx,
		cmdChan:       make(chan *client.Command),
		cmdHandlers:   map[string]CommandHandler{},
		LockTable:     locking.NewLockTable(),
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

	c := client.New(s.ctx, conn, s.cmdChan)

	c.ListenForCommands()

	s.LockTable.UnlockAllForClient(c)
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

func (s *Server) Write(c *client.Client, msg string) {
	_, _ = c.Write([]byte(msg + "\n"))
}
