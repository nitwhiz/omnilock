package server

import (
	"runtime"
	"time"
)

type Options struct {
	keepAlive     time.Duration
	acceptorCount int
	listenAddr    string
	clientTimeout time.Duration
}

type Option func(*Options)

func (o *Options) applyDefaults() {
	if o.keepAlive == 0 {
		WithKeepAlivePeriod(time.Second * 5)(o)
	}

	if o.acceptorCount == 0 {
		WithAcceptorCount(runtime.NumCPU())(o)
	}

	if o.listenAddr == "" {
		WithListenAddr("0.0.0.0:7194")(o)
	}

	if o.clientTimeout == 0 {
		WithClientTimeout(time.Second * 5)(o)
	}
}

func WithKeepAlivePeriod(p time.Duration) Option {
	return func(o *Options) {
		o.keepAlive = p
	}
}

func WithAcceptorCount(acceptorCount int) Option {
	return func(o *Options) {
		o.acceptorCount = acceptorCount
	}
}

func WithListenAddr(listenAddr string) Option {
	return func(o *Options) {
		o.listenAddr = listenAddr
	}
}

func WithClientTimeout(clientTimeout time.Duration) Option {
	return func(o *Options) {
		o.clientTimeout = clientTimeout
	}
}
