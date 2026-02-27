package internal

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
)

const (
	portsNum     int = 65535
	osMaxSockets int = 1150
)

type portStatus struct {
	num    int
	status bool
}

func WorkerPool(protocol string, target string) {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	inputs := make(chan int, 100)
	go func() {
		defer close(inputs)
		for i := 1; i <= portsNum; i++ {
			select {
			case <-ctx.Done():
				return
			case inputs <- i:
			}
		}
	}()

	results := make(chan portStatus, 100)
	var wg sync.WaitGroup
	wg.Add(osMaxSockets)

	for range osMaxSockets {
		go func() {
			defer wg.Done()

			for port := range inputs {
				res := ScanPort(ctx, protocol, target, port)
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
