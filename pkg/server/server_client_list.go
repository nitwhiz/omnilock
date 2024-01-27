package server

import (
	"sync"
)

type ClientList struct {
	autoIncrement uint64
	clients       map[uint64]*Client
	mu            *sync.RWMutex
}

func NewClientList() *ClientList {
	return &ClientList{
		autoIncrement: 0,
		clients:       map[uint64]*Client{},
		mu:            &sync.RWMutex{},
	}
}

func (l *ClientList) Add(c *Client) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.autoIncrement += 1

	c.ID = l.autoIncrement
	l.clients[l.autoIncrement] = c
}

func (l *ClientList) Remove(c *Client) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.clients, c.ID)
}
