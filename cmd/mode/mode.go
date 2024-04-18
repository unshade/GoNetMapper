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

		} else if mode == "server" {
			controller.ServerMode()
		} else {
			// normal mode
		}
	},
}
