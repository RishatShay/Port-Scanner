package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/RishatShay/Port-Scanner/internal/scanner"
)

func main() {
	host := flag.String("host", "scanme.nmap.org", "target host to scan")
	flag.Parse()

	protocol := "tcp"

	start := time.Now()
	scanner.WorkerPool(protocol, *host)
	fmt.Println(time.Since(start))
}
