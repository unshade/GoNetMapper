package internal

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
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

		buffer = buffer[:n]
		buffer = cleanBuffer(buffer)

		args := strings.Split(string(buffer), " ")
		commandName := args[0]

		cmdList := rootCmd.Commands()

		for _, cmd := range cmdList {
			if commandName == cmd.Name() {
				cmd.SetArgs(args[1:])

				pipeReader, pipeWriter, err := os.Pipe()
				if err != nil {
					return
				}

				cmd.SetOut(pipeWriter)
				cmd.SetErr(pipeWriter)

				stdout := os.Stdout
				stderr := os.Stderr

				os.Stdout = pipeWriter
				os.Stderr = pipeWriter

				cmd.Run(cmd, args[1:])

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
