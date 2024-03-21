package main

import (
	"fmt"
	"net"
)

func getGateways() ([]net.IPNet, error) {
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

func scanGateway(gateway net.IPNet) {
	mask := gateway.Mask
	network := gateway.IP.Mask(mask)
	broadcast := net.IP(make([]byte, len(network)))

	// Apply the bitwise NOT operator to the mask and OR it with the network address to get the broadcast address.
	for i := range network {
		broadcast[i] = network[i] | ^mask[i]
	}

	fmt.Printf("Network: %s\n", network.String())
	fmt.Printf("Broadcast: %s\n", broadcast.String())

	for ip := network.Mask(mask); !ip.Equal(broadcast); incIP(ip) {
		if !ip.Equal(network) && !ip.Equal(broadcast) {
			pingIp(ip)
		}
	}
}

func pingIp(ip net.IP) {
	fmt.Printf("Pinging %s\n", ip.String())

}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
