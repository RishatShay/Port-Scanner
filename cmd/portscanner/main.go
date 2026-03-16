package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/RishatShay/Port-Scanner/internal/scanner"
)

func main() {
	host := flag.String("host", "scanme.nmap.org", "target host to scan")
	protocol := flag.String("protocol", "tcp", "protocol to use: tcp or udp")
	portsFlag := flag.String("ports", "1-65535", "ports to scan, e.g. 22,80,443,8000-8100")
	timeout := flag.Duration("timeout", time.Second, "timeout per port")
	workers := flag.Int("workers", 500, "number of concurrent workers")
	flag.Parse()

	ports, err := scanner.ParsePorts(*portsFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid ports:", err)
		os.Exit(1)
	}

	cfg := scanner.Config{
		Host:     *host,
		Protocol: *protocol,
		Ports:    ports,
		Workers:  *workers,
		Timeout:  *timeout,
	}

	start := time.Now()
	scanner.WorkerPool(cfg)
	fmt.Println(time.Since(start))
}
