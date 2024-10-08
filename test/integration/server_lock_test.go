package integration

import (
	"bufio"
	"net"
	"sync"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	wg, serverAddr, cancel := startServer(t)

	defer func() {
		cancel()
		wg.Wait()
	}()

	conn := connect(t, serverAddr)

	defer func(conn *net.TCPConn) {
		_ = conn.Close()
	}(conn)

	_, err := conn.Write([]byte("lock test1\n"))

	if err != nil {
		t.Fatal(err)
	}

	err = conn.SetReadDeadline(time.Now().Add(time.Second))

	if err != nil {
		t.Fatal(err)
	}

	r := bufio.NewReader(conn)

	recv, err := r.ReadString('\n')

	if err != nil {
		t.Fatal(err)
	}

	recvEscaped := escapeString(recv)
	expected := "success\\x0A"

	if recvEscaped != expected {
		t.Fatalf("response mismatch, expected \"%s\", got \"%s\"", expected, recvEscaped)
	}
}

func TestLockWithTimeoutAndLockoutSameConnection(t *testing.T) {
	wg, serverAddr, cancel := startServer(t)

	defer func() {
		cancel()
		wg.Wait()
	}()

	conn := connect(t, serverAddr)

	defer func(conn *net.TCPConn) {
		_ = conn.Close()
	}(conn)

	// try 1 - success

	_, err := conn.Write([]byte("lock test1\n"))

	if err != nil {
		t.Fatal(err)
	}

	err = conn.SetReadDeadline(time.Now().Add(time.Second))

	if err != nil {
		t.Fatal(err)
	}

	r := bufio.NewReader(conn)

	recv, err := r.ReadString('\n')

	if err != nil {
		t.Fatal(err)
	}

	recvEscaped := escapeString(recv)
	expected := "success\\x0A"

	if recvEscaped != expected {
		t.Fatalf("response mismatch, expected \"%s\", got \"%s\"", expected, recvEscaped)
	}

	// try 2 - failure

	_, err = conn.Write([]byte("lock test1 200\n"))

	if err != nil {
		t.Fatal(err)
	}

	err = conn.SetReadDeadline(time.Now().Add(time.Second))

	if err != nil {
		t.Fatal(err)
	}

	recv, err = r.ReadString('\n')

	if err != nil {
		t.Fatal(err)
	}

	recvEscaped = escapeString(recv)
	expected = "failed\\x0A"

	if recvEscaped != expected {
		t.Fatalf("response mismatch, expected \"%s\", got \"%s\"", expected, recvEscaped)
	}
}

func TestLockWithTimeoutAndLockoutDifferentConnection(t *testing.T) {
	wg, serverAddr, cancel := startServer(t)

	defer func() {
		cancel()
		wg.Wait()
	}()

	// try 1 - success

	conn1 := connect(t, serverAddr)

	defer func(conn *net.TCPConn) {
		_ = conn.Close()
	}(conn1)

	_, err := conn1.Write([]byte("lock test1\n"))

	if err != nil {
		t.Fatal(err)
	}

	err = conn1.SetReadDeadline(time.Now().Add(time.Second))

	if err != nil {
		t.Fatal(err)
	}

	r1 := bufio.NewReader(conn1)

	recv, err := r1.ReadString('\n')

	if err != nil {
		t.Fatal(err)
	}

	recvEscaped := escapeString(recv)
	expected := "success\\x0A"

	if recvEscaped != expected {
		t.Fatalf("response mismatch, expected \"%s\", got \"%s\"", expected, recvEscaped)
	}

	// try 2 - failure

	conn2 := connect(t, serverAddr)

	defer func(conn *net.TCPConn) {
		_ = conn.Close()
	}(conn2)

	_, err = conn2.Write([]byte("lock test1 200\n"))

	if err != nil {
		t.Fatal(err)
	}

	err = conn2.SetReadDeadline(time.Now().Add(time.Second))

	if err != nil {
		t.Fatal(err)
	}

	r2 := bufio.NewReader(conn2)

	recv, err = r2.ReadString('\n')

	if err != nil {
		t.Fatal(err)
	}

	recvEscaped = escapeString(recv)
	expected = "failed\\x0A"

	if recvEscaped != expected {
		t.Fatalf("response mismatch, expected \"%s\", got \"%s\"", expected, recvEscaped)
	}
}

func TestLockWithTimeoutAndConcurrency(t *testing.T) {
	wg, serverAddr, cancel := startServer(t)

	defer func() {
		cancel()
		wg.Wait()
	}()

	wg2 := &sync.WaitGroup{}

	wg2.Add(1)

	go func() {
		conn1 := connect(t, serverAddr)

		defer func(conn *net.TCPConn) {
			_ = conn.Close()
			wg2.Done()
		}(conn1)

		_, err := conn1.Write([]byte("lock test1\n"))

		if err != nil {
			t.Error(err)
			return
		}

		r1 := bufio.NewReader(conn1)

		recv, err := r1.ReadString('\n')

		if err != nil {
			t.Error(err)
			return
		}

		recvEscaped := escapeString(recv)
		expected := "success\\x0A"

		if recvEscaped != expected {
			t.Errorf("response mismatch, expected \"%s\", got \"%s\"", expected, recvEscaped)
			return
		}

		time.Sleep(time.Millisecond * 250)
	}()

	wg2.Add(1)

	go func() {
		conn1 := connect(t, serverAddr)

		defer func(conn *net.TCPConn) {
			_ = conn.Close()
			wg2.Done()
		}(conn1)

		_, err := conn1.Write([]byte("lock test1 500\n"))

		if err != nil {
			t.Error(err)
			return
		}

		r1 := bufio.NewReader(conn1)

		recv, err := r1.ReadString('\n')

		if err != nil {
			t.Error(err)
			return
		}

		recvEscaped := escapeString(recv)
		expected := "success\\x0A"

		if recvEscaped != expected {
			t.Errorf("response mismatch, expected \"%s\", got \"%s\"", expected, recvEscaped)
			return
		}
	}()

	wg2.Wait()
}
