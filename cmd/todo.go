package cmd

import (
	"fmt"

	"github.com/loctherapy/dashboard/internal/todo"
	"github.com/spf13/cobra"
)

var pattern string

var todoCmd = &cobra.Command{
    Use:   "todo",
    Short: "List all the undone todos",
    Long:  `List all the undone todos in the current directory and nested directories`,
    Run: func(cmd *cobra.Command, args []string) {
        f := todo.NewFileFetcher(pattern)
        files, err := f.Fetch()
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        for _, file := range files {
            fmt.Println(file)
        }
    },
}

func init() {
    rootCmd.AddCommand(todoCmd)
    todoCmd.Flags().StringVarP(&pattern, "pattern", "p", `.*\.md$`, "Pattern to match files (default: *.md)")
}
