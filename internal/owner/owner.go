// Package owner finds the local process that holds an open TCP/UDP port.
//
// This only makes sense for ports on the machine the code runs on: there is
// no way to learn this over the network for a remote host, since the
// information lives in the target's own process table, not in the
// TCP/UDP handshake.
package owner

import "net"

// Owner describes the local process that holds an open socket.
type Owner struct {
	PID  int
	Name string
}

// Lookup tries to find the process that has the given local port open for
// protocol ("tcp" or "udp").
//
// The real lookup is OS-specific, see lookup_linux.go and lookup_windows.go.
// Everywhere else (lookup_fallback.go) it's simply not implemented. When the
// lookup isn't supported, fails, or the process can't be identified (e.g.
// it's owned by another user), ok is false and callers should fall back to
// reporting the port as open without an owner.
func Lookup(protocol string, port int) (Owner, bool) {
	return lookup(protocol, port)
}

// IsLocalHost reports whether host refers to the machine this code runs on.
// Lookup is only meaningful for such hosts.
func IsLocalHost(host string) bool {
	if host == "localhost" {
		return true
	}

	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}
	for _, ip := range ips {
		if ip.IsLoopback() {
			return true
		}
	}
	return false
}
