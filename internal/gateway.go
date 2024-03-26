package internal

import (
	"fmt"
	"net"
	"strings"
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

func ScanGatewayNetwork(gateway net.IPNet) string {
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

	n := 255

	channels := make([]chan string, n)

	for i := 0; i < n; i++ {
		channels[i] = make(chan string)
		go func(i int, ch chan string) {
			for ip := range ch {
				ip := ip
				// may be too fast
				go func() {
					duration, err := Ping(ip, 1*time.Second)
					if err == nil {
						fmt.Printf("Host %s is up: time=%v\n", ip, duration)
					}
				}()
			}
		}(i, channels[i])
	}

	c := 0
	for ip := network.Mask(mask); !ip.Equal(broadcast); IncIP(ip) {
		if !ip.Equal(network) && !ip.Equal(broadcast) {
			channels[c%n] <- ip.String()
			c++
		}
	}

	for i := 0; i < n; i++ {
		close(channels[i])
	}

	return builder.String()
}
