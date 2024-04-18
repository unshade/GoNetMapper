package mode

import (
	"github.com/spf13/cobra"
	"main/internal/controller"
	"sync"
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
			waitGroup := sync.WaitGroup{}
			waitGroup.Add(1)
			controller.ServerMode()
			waitGroup.Wait()
		} else {
			// normal mode
		}
	},
}
