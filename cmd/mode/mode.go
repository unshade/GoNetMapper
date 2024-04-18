package mode

import (
	"github.com/spf13/cobra"
	"main/internal/controller"
)

var ExecMode = &cobra.Command{
	Use:   "exec-mode",
	Short: "Execute the app in client or server mode",
	Long:  "Execute the app in client or server mode",
	Run: func(cmd *cobra.Command, args []string) {
		mode := args[0]
		if mode == "client" {
			ip := ""
			if len(args) == 2 {
				ip = args[1]
			} else {
				ip = "127.0.0.1"
			}
			controller.ClientMode(ip)
		} else if mode == "server" {
			controller.ServerMode()
			select {}
		} else {
			// normal mode
		}
	},
}
