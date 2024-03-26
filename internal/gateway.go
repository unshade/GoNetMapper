package internal

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
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

func ScanGateway(gateway net.IPNet) {
	mask := gateway.Mask
	network := gateway.IP.Mask(mask)
	broadcast := net.IP(make([]byte, len(network)))

	// Apply the bitwise NOT operator to the mask and OR it with the network address to get the broadcast address.
	for i := range network {
		broadcast[i] = network[i] | ^mask[i]
	}

	fmt.Printf("Network: %s\n", network.String())
	fmt.Printf("Broadcast: %s\n", broadcast.String())

	// Prepare the goroutines
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
	for ip := network.Mask(mask); !ip.Equal(broadcast); incIP(ip) {
		if !ip.Equal(network) && !ip.Equal(broadcast) {
			channels[c%n] <- ip.String()
			c++
		}
	}

	for i := 0; i < n; i++ {
		close(channels[i])
	}
}

func Ping(ipString string, timeout time.Duration) (time.Duration, error) {
	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return 0, err
	}
	defer c.Close()
	// Generate an Echo message
	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("Hello, are you there!"),
		},
	}
	wb, err := msg.Marshal(nil)
	if err != nil {
		return 0, err
	}
	// Send, note that here it must be a UDP address
	start := time.Now()
	if _, err := c.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP(ipString)}); err != nil {
		return 0, err
	}
	// Read the reply package
	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return 0, err
	}
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return 0, err
	}
	duration := time.Since(start)

	// The reply packet is an ICMP message, parsed first
	msg, err = icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return 0, err
	}

	// Get results
	if msg.Type == ipv4.ICMPTypeEchoReply {
		echoReply, ok := msg.Body.(*icmp.Echo) // The message body is of type Echo
		if !ok {
			return 0, fmt.Errorf("invalid ICMP Echo Reply message")
		}
		if peer.(*net.UDPAddr).IP.String() == ipString && echoReply.Seq == 1 && echoReply.ID == os.Getpid()&0xffff {
			return duration, nil
		}
	}
	return 0, fmt.Errorf("unexpected ICMP message type: %v", msg.Type)
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
