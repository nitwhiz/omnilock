package integration

import (
	"context"
	"fmt"
	"github.com/nitwhiz/omnilock/v2/pkg/server"
	"net"
	"strings"
	"sync"
	"testing"
)

func connect(t *testing.T, tcpServer *net.TCPAddr) *net.TCPConn {
	conn, err := net.DialTCP("tcp", nil, tcpServer)

	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func startServer(t *testing.T) (*sync.WaitGroup, *net.TCPAddr, context.CancelFunc) {
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	s, err := server.New(ctx, server.WithListenAddr("localhost:3000"))

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		wg.Add(1)
		defer wg.Done()

		s.Run()
	}()

	tcpServer, err := net.ResolveTCPAddr("tcp", "localhost:3000")

	return wg, tcpServer, cancel
}

func escapeString(input string) string {
	var escaped strings.Builder

	for _, runeValue := range input {
		if runeValue >= ' ' && runeValue <= '~' {
			escaped.WriteRune(runeValue)
		} else {
			escaped.WriteString(fmt.Sprintf("\\x%02X", runeValue))
		}
	}

	return escaped.String()
}
