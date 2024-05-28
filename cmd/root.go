package cmd

import (
	"fmt"
	"os"

	"github.com/loctherapy/dashboard/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Dashboard is a CLI tool for managing your application",
	Long:  `A longer description of Dashboard with usage examples and details.`,
	Run: func(cmd *cobra.Command, args []string) {
		uiFactory := ui.NewUIFactory()
		uiFactory.Run()
	},
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Define flags and configuration settings here if needed
}
