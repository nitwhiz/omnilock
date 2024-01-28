package integration

import (
	"context"
	"fmt"
	"github.com/nitwhiz/omnilock/pkg/server"
	"net"
	"strings"
	"testing"
)

func connect(t *testing.T, tcpServer *net.TCPAddr) *net.TCPConn {
	conn, err := net.DialTCP("tcp", nil, tcpServer)

	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func startServer(t *testing.T) (*net.TCPAddr, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	s, err := server.New(ctx, server.WithListenAddr("localhost:3000"))

	if err != nil {
		t.Fatal(err)
	}

	go s.Accept()

	tcpServer, err := net.ResolveTCPAddr("tcp", "localhost:3000")

	return tcpServer, cancel
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
