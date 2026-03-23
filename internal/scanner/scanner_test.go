package scanner

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestScanPort_Open(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	host, port := hostAndPort(t, listener.Addr())

	if !ScanPort(context.Background(), "tcp", host, port, time.Second) {
		t.Fatalf("expected port %d on %s to be open", port, host)
	}
}

func TestScanPort_Closed(t *testing.T) {
	// Reserve a port and close it right away, so nothing listens on it.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to reserve a port: %v", err)
	}
	host, port := hostAndPort(t, listener.Addr())
	listener.Close()

	if ScanPort(context.Background(), "tcp", host, port, 200*time.Millisecond) {
		t.Fatalf("expected port %d on %s to be closed", port, host)
	}
}

func hostAndPort(t *testing.T, addr net.Addr) (string, int) {
	t.Helper()

	host, portStr, err := net.SplitHostPort(addr.String())
	if err != nil {
		t.Fatalf("failed to split address %q: %v", addr, err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse port from %q: %v", portStr, err)
	}

	return host, port
}
