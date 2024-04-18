package internal

import (
	"fmt"
	"net"
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

func ScanGatewayNetwork(gateway net.IPNet, progression *float64) {
	*progression = 0

	mask := gateway.Mask
	network := gateway.IP.Mask(mask)
	broadcast := net.IP(make([]byte, len(network)))

	// Apply the bitwise NOT operator to the mask and OR it with the network address to get the broadcast address.
	for i := range network {
		broadcast[i] = network[i] | ^mask[i]
	}

	fmt.Printf("Network: %s\n", network.String())
	fmt.Printf("Broadcast: %s\n", broadcast.String())

	mut := sync.Mutex{}
	wg := sync.WaitGroup{}

	done := 0
	_, bits := mask.Size()
	total := 1 << uint(32-bits)

	list := make([]string, 0, total)

	c := 0
	for ip := network.Mask(mask); !ip.Equal(broadcast); ip = NextIP(ip, 1) {
		if !ip.Equal(network) && !ip.Equal(broadcast) {
			list = append(list, ip.String())
			c++
		}
	}

	for _, ip := range list {
		ip := ip
		// may be too fast
		wg.Add(1)
		go func() {
			defer wg.Done()
			duration, err := Ping(ip, 1*time.Second)

			mut.Lock()
			done++
			*progression = float64(done) / float64(total)
			mut.Unlock()

			if err == nil {
				fmt.Printf("Host %s is up: time=%v\n", ip, duration)
			}
		}()
	}

	wg.Wait()
}
