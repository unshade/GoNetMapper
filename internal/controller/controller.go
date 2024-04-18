package controller

import (
	"fmt"
	"io"
	"log"
	"main/cmd/scan_gateway"
	"main/cmd/scan_ports"
	"net"
	"os"
	"strings"
)

func ServerMode() {

	listen, err := net.Listen("tcp", "127.0.0.1:6666")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println("Daemon listening on port 6666")

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			go handleConnection(conn)
		}
	}()
}

func handleConnection(conn net.Conn) {

	conn.Write([]byte("You are connected to Radar Daemon on TCP port 6666\n"))
	for {
		buffer := make([]byte, 1024)

		n, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed")
			} else {
				fmt.Println("Error reading from connection", err)
			}
			return
		}

		buffer = buffer[:n]
		buffer = cleanBuffer(buffer)

		args := strings.Split(string(buffer), " ")
		commandName := args[0]

		if commandName == "exit" {
			conn.Close()
			return
		}

		pipeReader, pipeWriter, err := os.Pipe()
		if err != nil {
			return
		}

		stdout := os.Stdout
		stderr := os.Stderr

		os.Stdout = pipeWriter
		os.Stderr = pipeWriter

		switch args[0] {
		case "scan-ports":
			fmt.Println("Running scan-ports")
			scan_ports.ScanPortsCommand.Run(scan_ports.ScanPortsCommand, args[1:])
		case "scan-gateways":
			fmt.Println("Running scan-gateways")
			scan_gateway.ScanGatewayCommand.Run(scan_gateway.ScanGatewayCommand, args[1:])
		}

		//cmd.Run(cmd, args[1:])

		pipeWriter.Close()

		_, err = io.Copy(conn, pipeReader)
		if err != nil {
			fmt.Println("Error copying to connection", err)
			return
		}
		pipeReader.Close()

		os.Stdout = stdout
		os.Stderr = stderr
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
