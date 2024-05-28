package cmd

import (
	"fmt"

	"github.com/loctherapy/dashboard/internal/todo"
	"github.com/spf13/cobra"
)

var pattern string

var todoCmd = &cobra.Command{
    Use:   "todo",
    Short: "List all unchecked todos in markdown files",
    Long:  `Fetch all markdown files and list all unchecked todos within them.`,
    Run: func(cmd *cobra.Command, args []string) {
        facade, err := todo.NewToDoFacade(`.*\.md$`, todo.Console)

        if err != nil {
            fmt.Println("Error printing todos:", err)
        }

        todosString, err := facade.Print()

        if err != nil {
            fmt.Println("Error printing todos:", err)
        }

        fmt.Println(todosString)
    },
}


func init() {
    rootCmd.AddCommand(todoCmd)
    todoCmd.Flags().StringVarP(&pattern, "pattern", "p", `.*\.md$`, "Pattern to match files (default: *.md)")
}
