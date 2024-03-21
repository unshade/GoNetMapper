package main

import (
	"fmt"
	"log"
)

func main() {

	gateways, err := getGateways()

	if err != nil {
		log.Fatal(err)
	}

	for _, gateway := range gateways {
		fmt.Println(gateway)
		scanGateway(gateway)
	}
}
