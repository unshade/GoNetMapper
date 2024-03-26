package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/internal"
)

func init() {
	rootCmd.AddCommand(scanPortsCommand)
}

var scanPortsCommand = &cobra.Command{
	Use:   "scan-ports",
	Short: "Scan ports of an ip address",
	Long:  "Scan ports of an ip address",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires an ip address")
		}
		ip := args[0]
		if !internal.IsValidIP(ip) {
			return fmt.Errorf("invalid ip address")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		internal.TcpScan(ip, 1, 10000)
	},
}
