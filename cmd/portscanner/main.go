package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/RishatShay/Port-Scanner/internal/scanner"
)

func main() {
	host := flag.String("host", "scanme.nmap.org", "target host to scan")
	protocol := flag.String("protocol", "tcp", "protocol to use: tcp or udp")
	timeout := flag.Duration("timeout", time.Second, "timeout per port")
	flag.Parse()

	cfg := scanner.Config{
		Host:     *host,
		Protocol: *protocol,
		Timeout:  *timeout,
	}

	start := time.Now()
	scanner.WorkerPool(cfg)
	fmt.Println(time.Since(start))
}
