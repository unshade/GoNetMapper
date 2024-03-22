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
				ip := net.ParseIP(ip)
				err := Ping(ip, 1*time.Second)
				if err == nil {
					fmt.Printf("Host %s is up\n", ip.String())
				}
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

func Ping(ip net.IP, timeout time.Duration) error {
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return err
	}
	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte(""),
		},
	}
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return err
	}

	// Write the message to the listening connection
	if _, err := conn.WriteTo(msgBytes, &net.UDPAddr{IP: net.ParseIP(ip.String())}); err != nil {
		return err
	}

	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}
	reply := make([]byte, 1500)
	n, _, err := conn.ReadFrom(reply)

	if err != nil {
		return err
	}
	parsedReply, err := icmp.ParseMessage(1, reply[:n])

	if err != nil {
		return err
	}

	if parsedReply.Type == ipv4.ICMPTypeEchoReply {
		return nil
	}

	return fmt.Errorf("reply from %s is not an echo reply", ip.String())
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
