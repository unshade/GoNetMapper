package main

import (
	"fmt"
	"log"
	"main/internal"
)

func main() {

	gateways, err := internal.GetGateways()

	if err != nil {
		log.Fatal(err)
	}

	for _, gateway := range gateways {
		fmt.Println(gateway)
		internal.ScanGateway(gateway)
		internal.TcpScan(gateway.IP.To4().String(), 1, 10000)
	}
}
