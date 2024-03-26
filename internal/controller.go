package internal

import (
	"fmt"
	"log"
	"net"
	"os"
)

func ServerMode() {

	listen, err := net.Listen("tcp", "127.0.0.1:6666")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer listen.Close()

	log.Println("Daemon listening on port 6666")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {

	conn.Write([]byte("You are connected to Radar Daemon on TCP port 6666\n"))
	for {
		buffer := make([]byte, 1024)

		_, err := conn.Read(buffer)

		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		fmt.Println(len(buffer))
		fmt.Println(string(buffer))
		switch string(cleanBuffer(buffer)) {
		case "ping":
			conn.Write([]byte("pong\n"))
		case "status":
			conn.Write([]byte("daemon is running\n"))
		case "scan-gateway":
			gateway, _ := GetGateways()
			res := ScanGatewayNetwork(gateway[0])
			conn.Write([]byte(res))
		case "kill":
			conn.Write([]byte("killing daemon\n"))
			conn.Close()
		default:
			conn.Write([]byte("unknown command\n"))
		}
	}
}

func cleanBuffer(buffer []byte) []byte {

	for i := 0; i < len(buffer); i++ {
		if buffer[i] == '\n' || buffer[i] == '\r' {
			buffer = buffer[:i]
		}
	}

	return buffer
}
