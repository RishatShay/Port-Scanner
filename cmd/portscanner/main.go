package main

import (
	"fmt"
	"time"

	"github.com/RishatShay/Port-Scanner/internal"
)

func main() {
	start := time.Now()
	target := "scanme.nmap.org"
	protocol := "tcp"

	internal.WorkerPool(protocol, target)
	fmt.Println(time.Since(start))
}
