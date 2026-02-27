package internal

import (
	"context"
	"fmt"
	"net"
	"time"
)

func ScanPort(ctx context.Context, protocol string, target string, port int) bool {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
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
