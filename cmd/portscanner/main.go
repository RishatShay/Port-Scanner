package main

import (
	"fmt"
	"time"

	"github.com/RishatShay/Port-Scanner/internal/scanner"
)

func main() {
	start := time.Now()
	target := "scanme.nmap.org"
	protocol := "tcp"

	scanner.WorkerPool(protocol, target)
	fmt.Println(time.Since(start))
}
