package cmd

import (
	"fmt"
	"os"

	"github.com/loctherapy/dashboard/internal/controller"
	"github.com/loctherapy/dashboard/internal/repository"
	"github.com/loctherapy/dashboard/internal/service"
	"github.com/loctherapy/dashboard/internal/view"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Dashboard is a CLI tool for managing your application",
	Long:  `A longer description of Dashboard with usage examples and details.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Model
		fileFetcher, err := repository.NewFileFetcher(`.*\.md$`)
		if err != nil { 
			fmt.Println(err)
			os.Exit(1)
		}
		todoRepository := repository.NewToDoRepository(fileFetcher)
		todoService := service.NewToDoService(todoRepository)

		// View
		factory := view.ToDoPrinterFactory{}
		todoPrinter, err := factory.CreatePrinter(view.TView)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		
		todoView := view.NewView(todoPrinter)

		// Controller
		todoController := controller.NewToDoController(todoService, todoView)	
		
		todoController.GetTodos()
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
