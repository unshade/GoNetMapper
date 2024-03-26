package internal

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func scanPort(ip string, port int, results chan<- int) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	results <- port
}

func TcpScan(ip string, startPort, endPort int) {
	fmt.Println("Scanning ports on", ip)

	results := make(chan int, endPort-startPort+1)

	var wg sync.WaitGroup
	wg.Add(endPort - startPort + 1)

	for port := startPort; port <= endPort; port++ {
		go func(p int) {
			defer wg.Done()
			scanPort(ip, p, results)
		}(port)
	}

	wg.Wait()
	close(results)

	for port := range results {
		fmt.Printf("Port %d is open\n", port)
	}
}
