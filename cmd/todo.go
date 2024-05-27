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
        // Fetch .md files
        f := todo.NewFileFetcher(`.*\.md$`)
        files, err := f.Fetch()
        if err != nil {
            fmt.Println("Error fetching files:", err)
            return
        }

        // Build todos
        builder := todo.NewToDoBuilder()
        todos, err := builder.Build(files)
        if err != nil {
            fmt.Println("Error building todos:", err)
            return
        }

        // Print results
        for _, fileToDos := range todos {
            fmt.Printf("File: %s\n", fileToDos.FilePath)
            for _, todo := range fileToDos.ToDos {
                fmt.Println(todo.Line)
            }
            fmt.Println()
        }
    },
}


func init() {
    rootCmd.AddCommand(todoCmd)
    todoCmd.Flags().StringVarP(&pattern, "pattern", "p", `.*\.md$`, "Pattern to match files (default: *.md)")
}
