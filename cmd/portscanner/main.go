package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RishatShay/Port-Scanner/internal/owner"
	"github.com/RishatShay/Port-Scanner/internal/scanner"
)

func main() {
	host := flag.String("host", "scanme.nmap.org", "target host to scan")
	protocol := flag.String("protocol", "tcp", "protocol to use: tcp or udp")
	portsFlag := flag.String("ports", "1-65535", "ports to scan, e.g. 22,80,443,8000-8100")
	timeout := flag.Duration("timeout", time.Second, "timeout per port")
	workers := flag.Int("workers", 500, "number of concurrent workers")
	flag.Parse()

	if *protocol != "tcp" && *protocol != "udp" {
		fmt.Fprintln(os.Stderr, "protocol must be either tcp or udp")
		os.Exit(1)
	}

	if *host == "" {
		fmt.Fprintln(os.Stderr, "host must not be empty")
		os.Exit(1)
	}

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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	local := owner.IsLocalHost(*host)

	start := time.Now()
	results := scanner.Run(ctx, cfg)

	for _, r := range results {
		if local {
			if o, ok := owner.Lookup(*protocol, r.Port); ok {
				fmt.Printf("%d/%s is open (pid %d, %s)\n", r.Port, *protocol, o.PID, o.Name)
				continue
			}
		}
		fmt.Printf("%d/%s is open\n", r.Port, *protocol)
	}
	fmt.Printf("scanned %d ports in %s, found %d open\n", len(ports), time.Since(start), len(results))
}
