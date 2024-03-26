package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "radar",
	Short: "Radar is an easy network scanner and monitoring tool.",
	Long:  "Radar is an easy network scanner and monitoring tool.",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := tea.NewProgram(newModel()).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
