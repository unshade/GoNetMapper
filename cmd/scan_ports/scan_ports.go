package scan_ports

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/cmd"
	"main/internal"
	"strconv"
)

func init() {
	cmd.RootCmd.AddCommand(scanPortsCommand)
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

		// Check start and end port
		if len(args) >= 2 {
			startPort := args[1]
			if !internal.IsValidPort(startPort) {
				return fmt.Errorf("invalid start port")
			}

			if len(args) == 3 {
				endPort := args[2]
				if !internal.IsValidPort(endPort) {
					return fmt.Errorf("invalid end port")
				}

				if startPort > endPort {
					return fmt.Errorf("start port must be less than end port")
				}
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		startPort := 1
		endPort := 10000
		if len(args) >= 2 {
			startPort, _ = strconv.Atoi(args[1])
			if len(args) == 3 {
				endPort, _ = strconv.Atoi(args[2])
			}
		}
		internal.TcpScan(ip, startPort, endPort)
	},
}
