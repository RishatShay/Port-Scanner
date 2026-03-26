// Package scanner implements a concurrent TCP/UDP port scanner.
//
// A scan is described with a Config and executed with Run, which fans
// work out to a pool of goroutines and returns the ports that answered.
package scanner
