package scanner

import (
	"context"
	"fmt"
	"net"
	"time"
)

// ScanPort tries to open a connection to a single port and reports whether
// it is open. It gives up after timeout.
func ScanPort(ctx context.Context, protocol string, target string, port int, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	address := fmt.Sprintf("%s:%d", target, port)
	var d net.Dialer
	conn, err := d.DialContext(ctx, protocol, address)

	if err != nil {
		return false
	}
	conn.Close()
	return true
}
