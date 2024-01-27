package server

import (
	"time"
)

type Option func(*Server)

func WithKeepAlivePeriod(p time.Duration) Option {
	return func(s *Server) {
		s.keepAlive = p
	}
}

func WithAcceptorCount(acceptorCount int) Option {
	return func(s *Server) {
		s.acceptorCount = acceptorCount
	}
}

func WithListenAddr(listenAddr string) Option {
	return func(s *Server) {
		s.listenAddr = listenAddr
	}
}
