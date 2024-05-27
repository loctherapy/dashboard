package cmd

import (
	"fmt"

	"github.com/loctherapy/dashboard/internal/fetcher"
	"github.com/spf13/cobra"
)

var pattern string

var fetchCmd = &cobra.Command{
    Use:   "fetch",
    Short: "Fetch and list all files matching a pattern",
    Long:  `Fetch and list all files in the current directory and nested directories matching the given pattern.`,
    Run: func(cmd *cobra.Command, args []string) {
        f := fetcher.NewFileFetcher(pattern)
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
    rootCmd.AddCommand(fetchCmd)
    fetchCmd.Flags().StringVarP(&pattern, "pattern", "p", `.*\.md$`, "Pattern to match files (default: *.md)")
}
