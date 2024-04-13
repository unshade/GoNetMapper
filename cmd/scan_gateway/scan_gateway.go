package scan_gateway

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/internal"
)

var ScanGatewayCommand = &cobra.Command{
	Use:   "scan-gateways",
	Short: "Scan the gateways of the current network",
	Long:  "Scan the gateways of the current network",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scanning gateways...")
		gateways, err := internal.GetGateways()

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, gateway := range gateways {
			fmt.Println(gateway)
			internal.ScanGatewayNetwork(gateway)
		}
	},
}
