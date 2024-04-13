package internal

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"strings"
)

func ServerMode(rootCmd *cobra.Command) {

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

		go handleConnection(conn, rootCmd)
	}

}

func handleConnection(conn net.Conn, rootCmd *cobra.Command) {

	conn.Write([]byte("You are connected to Radar Daemon on TCP port 6666\n"))
	for {
		buffer := make([]byte, 1024)

		n, err := conn.Read(buffer)

		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		buffer = cleanBuffer(buffer)
		buffer = buffer[:n]

		fmt.Println("Received buffer", string(buffer))
		args := strings.Split(string(buffer), " ")
		commandName := args[0]
		fmt.Println("Executing command", commandName, "with args", args[1:])

		cmdList := rootCmd.Commands()

		for _, cmd := range cmdList {
			if cmd.Name() == commandName {
				cmd.SetArgs(args[1:])
				fmt.Println("Executing command", cmd.Name())
				go cmd.Run(cmd, args[1:])
			}

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
