package cmd

import (
	"fmt"
	"sort"

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

        // Group by context
        contextMap := make(map[string][]todo.FileToDos)
        for _, fileToDos := range todos {
            contextMap[fileToDos.Context] = append(contextMap[fileToDos.Context], fileToDos)
        }

        // Print results
        for context, files := range contextMap {
            fmt.Printf("Context: %s\n", context)

            // Sort files by gravity
            sort.Slice(files, func(i, j int) bool {
                return files[i].Gravity > files[j].Gravity
            })

            for _, fileToDos := range files {
                fmt.Printf("  File: %s (Gravity: %d)\n", fileToDos.FilePath, fileToDos.Gravity)
                for _, todo := range fileToDos.ToDos {
                    fmt.Printf("    %s\n", todo.Line)
                }
            }
            fmt.Println()
        }
    },
}


func init() {
    rootCmd.AddCommand(todoCmd)
    todoCmd.Flags().StringVarP(&pattern, "pattern", "p", `.*\.md$`, "Pattern to match files (default: *.md)")
}
