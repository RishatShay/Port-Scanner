package scanner

import (
	"context"
	"sort"
	"sync"
	"time"
)

// defaultWorkers is used when Config.Workers is not set.
const defaultWorkers int = 500

// Config holds everything the scanner needs to know about a scan.
type Config struct {
	Host     string
	Protocol string
	Ports    []int
	Workers  int
	Timeout  time.Duration
}

// Result is the outcome of scanning a single port.
type Result struct {
	Port int
	Open bool
}

// Run scans all ports from cfg concurrently and returns the open ones,
// sorted by port number. It stops early if ctx is cancelled.
func Run(ctx context.Context, cfg Config) []Result {
	workers := cfg.Workers
	if workers <= 0 {
		workers = defaultWorkers
	}

	inputs := make(chan int, 100)
	go func() {
		defer close(inputs)
		for _, port := range cfg.Ports {
			select {
			case <-ctx.Done():
				return
			case inputs <- port:
			}
		}
	}()

	results := make(chan Result, 100)
	var wg sync.WaitGroup
	wg.Add(workers)

	for range workers {
		go func() {
			defer wg.Done()

			for port := range inputs {
				open := ScanPort(ctx, cfg.Protocol, cfg.Host, port, cfg.Timeout)
				results <- Result{Port: port, Open: open}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var open []Result
	for r := range results {
		if r.Open {
			open = append(open, r)
		}
	}

	sort.Slice(open, func(i, j int) bool { return open[i].Port < open[j].Port })
	return open
}
