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
	conn, err := net.DialTimeout("tcp", target, 3*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func TcpScan(ip string, startPort, endPort int, progression *float64) {
	fmt.Println("Scanning ports on", ip)
	*progression = 0

	var results sync.Map
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, maxGoroutines)
	mutex := &sync.Mutex{}

	done := 0

	for port := startPort; port <= endPort; port++ {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// Update progression
			defer func() {
				mutex.Lock()
				done++
				*progression = float64(done) / float64(endPort-startPort+1)
				mutex.Unlock()
			}()

			if isOpen := scanPort(ip, p); isOpen {
				results.Store(p, struct{}{})
				fmt.Printf("Port %d is open\n", p)
			}
		}(port)
	}

	wg.Wait()

	results.Range(func(key, value interface{}) bool {
		/*port := key.(int)
		fmt.Printf("Port %d is open\n", port)*/
		return true
	})
}
