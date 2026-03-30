//go:build linux

package owner

import (
	"net"
	"os"
	"strconv"
	"testing"
)

func TestLookup_Linux(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer listener.Close()

	port := portOf(t, listener.Addr())

	owner, ok := Lookup("tcp", port)
	if !ok {
		t.Fatalf("expected to find an owner for port %d", port)
	}
	if owner.PID != os.Getpid() {
		t.Errorf("owner.PID = %d, want %d (this test process)", owner.PID, os.Getpid())
	}
	if owner.Name == "" {
		t.Error("owner.Name is empty")
	}
}

func TestLookup_ClosedPort(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to reserve a port: %v", err)
	}
	port := portOf(t, listener.Addr())
	listener.Close()

	if _, ok := Lookup("tcp", port); ok {
		t.Errorf("expected no owner for closed port %d", port)
	}
}

func portOf(t *testing.T, addr net.Addr) int {
	t.Helper()

	_, portStr, err := net.SplitHostPort(addr.String())
	if err != nil {
		t.Fatalf("failed to split address %q: %v", addr, err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse port from %q: %v", portStr, err)
	}

	return port
}
