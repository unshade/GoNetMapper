package internal

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func GetGateways() ([]net.IPNet, error) {
	var gateways []net.IPNet

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					gateways = append(gateways, *ipnet)
				}
			}
		}
	}

	return gateways, nil
}

func ScanGatewayNetwork(gateway net.IPNet, progression *float64) string {
	*progression = 0

	mask := gateway.Mask
	network := gateway.IP.Mask(mask)
	broadcast := net.IP(make([]byte, len(network)))

	builder := strings.Builder{}

	// Apply the bitwise NOT operator to the mask and OR it with the network address to get the broadcast address.
	for i := range network {
		broadcast[i] = network[i] | ^mask[i]
	}

	fmt.Printf("Network: %s\n", network.String())
	fmt.Printf("Broadcast: %s\n", broadcast.String())
	builder.WriteString(fmt.Sprintf("Network: %s\n", network.String()))
	builder.WriteString(fmt.Sprintf("Broadcast: %s\n", broadcast.String()))
	builder.WriteString("Mask: " + mask.String() + "\n")

	var mut sync.Mutex
	var wg sync.WaitGroup

	done := 0
	_, bits := mask.Size()
	total := 1 << uint(32-bits)
	ch := make(chan string)

	go func() {
		for ip := range ch {
			ip := ip
			// may be too fast
			go func() {
				wg.Add(1)
				defer wg.Done()
				duration, err := Ping(ip, 1*time.Second)

				mut.Lock()
				done++
				*progression = float64(done) / float64(total)
				mut.Unlock()

				if err == nil {
					fmt.Printf("Host %s is up: time=%v\n", ip, duration)
					mut.Lock()
					builder.WriteString(fmt.Sprintf("Host %s is up: time=%v\n", ip, duration))
					mut.Unlock()
				}
			}()
		}
	}()

	c := 0
	for ip := network.Mask(mask); !ip.Equal(broadcast); IncIP(ip) {
		if !ip.Equal(network) && !ip.Equal(broadcast) {
			ch <- ip.String()
			c++
		}
	}
	close(ch)

	wg.Wait()

	return builder.String()
}
