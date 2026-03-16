package scanner

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
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

type portStatus struct {
	num    int
	status bool
}

func WorkerPool(cfg Config) {
	workers := cfg.Workers
	if workers <= 0 {
		workers = defaultWorkers
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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

	results := make(chan portStatus, 100)
	var wg sync.WaitGroup
	wg.Add(workers)

	for range workers {
		go func() {
			defer wg.Done()

			for port := range inputs {
				res := ScanPort(ctx, cfg.Protocol, cfg.Host, port, cfg.Timeout)
				results <- portStatus{num: port, status: res}
			}

		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for v := range results {
		if v.status {
			fmt.Printf("Port %d is open\n", v.num)
		}
	}
}
