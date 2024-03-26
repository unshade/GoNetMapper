package internal

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const maxGoroutines = 1000

func scanPort(ip string, port int) bool {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func TcpScan(ip string, startPort, endPort int) {
	fmt.Println("Scanning ports on", ip)

	var results sync.Map
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, maxGoroutines)

	for port := startPort; port <= endPort; port++ {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			defer func() { <-semaphore }()
			if isOpen := scanPort(ip, p); isOpen {
				results.Store(p, struct{}{})
			}
		}(port)
	}

	wg.Wait()

	results.Range(func(key, value interface{}) bool {
		port := key.(int)
		fmt.Printf("Port %d is open\n", port)
		return true
	})
}
