package internal

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

func IsValidIP(ip string) bool {
	regex := `^(\d{1,3}\.){3}\d{1,3}$`
	return regexp.MustCompile(regex).MatchString(ip)
}

func IsValidPort(port string) bool {
	i, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return i > 0 && i < 65536
}

func IncIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func NextIP(ip net.IP, inc uint) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
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
		if peer.(*net.UDPAddr).IP.String() == ipString && echoReply.Seq == 1 /*&& echoReply.ID == os.Getpid()&0xffff*/ {
			return duration, nil
		}
	}
	return 0, fmt.Errorf("unexpected ICMP message type: %v", msg.Type)
}
