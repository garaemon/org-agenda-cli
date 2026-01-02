package cmd

import (
	"fmt"
	"os"

	"github.com/garaemon/org-agenda-cli/pkg/config"
	"github.com/garaemon/org-agenda-cli/pkg/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	todoStatus   string
	todoTag      string
	todoFile     string
	todoSchedule string
	todoDeadline string
	todoTags     string
)

// todoCmd represents the todo command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manages TODO items",
	Long:  `Manages TODO items.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior: list
		todoListCmd.Run(cmd, args)
	},
}

var todoListCmd = &cobra.Command{
	Use:   "list",
	Short: "Display a list of TODO items",
	Run: func(cmd *cobra.Command, args []string) {
		paths := viper.GetStringSlice("org_files")
		if len(paths) == 0 {
			// Fallback for testing if no config exists
			if _, err := os.Stat("sample.org"); err == nil {
				paths = []string{"sample.org"}
			} else {
				fmt.Println("No org files configured. Use 'org-agenda config add-path <path>' to add one.")
				return
			}
		}

		orgFiles := config.ResolveOrgFiles(paths)
		for _, file := range orgFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file, err)
				continue
			}

			items := parser.ParseString(string(content), file)
			for _, item := range items {
				if todoStatus != "" {
					if item.Status != todoStatus {
						continue
					}
				} else {
					if item.Status == "" {
						continue
					}
				}
				// Basic filtering by tag (simple implementation)
				if todoTag != "" {
					found := false
					for _, t := range item.Tags {
						if t == todoTag {
							found = true
							break
						}
					}
					if !found {
						continue
					}
				}

				statusStr := item.Status
				if statusStr == "" {
					statusStr = "NONE"
				}
				fmt.Printf("[%s] %s (%s:%d)\n", statusStr, item.Title, item.FilePath, item.LineNumber)
			}
		}
	},
}

var todoAddCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new TODO item",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		targetFile := todoFile
		if targetFile == "" {
			targetFile = viper.GetString("default_file")
		}
		if targetFile == "" {
			orgFiles := viper.GetStringSlice("org_files")
			if len(orgFiles) > 0 {
				targetFile = orgFiles[0]
			}
		}

		if targetFile == "" {
			fmt.Println("Error: No target file specified and no default file configured.")
			return
		}

		f, err := os.OpenFile(targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Printf("Error closing file: %v\n", err)
			}
		}()

		content := fmt.Sprintf("* TODO %s", title)
		if todoTags != "" {
			content += fmt.Sprintf(" :%s:", todoTags)
		}
		content += "\n"

		if todoSchedule != "" {
			content += fmt.Sprintf("SCHEDULED: <%s>\n", todoSchedule)
		}
		if todoDeadline != "" {
			content += fmt.Sprintf("DEADLINE: <%s>\n", todoDeadline)
		}

		if _, err := f.WriteString(content); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}

		fmt.Printf("Added task: %s to %s\n", title, targetFile)
	},
}

var todoDoneCmd = &cobra.Command{
	Use:   "done [id|index]",
	Short: "Mark a task as DONE",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("todo done called for: %s\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)

	todoCmd.AddCommand(todoListCmd)
	todoCmd.AddCommand(todoAddCmd)
	todoCmd.AddCommand(todoDoneCmd)

	todoListCmd.Flags().StringVar(&todoStatus, "status", "", "Filter by status (TODO|WAITING|DONE)")
	todoListCmd.Flags().StringVar(&todoTag, "tag", "", "Filter by tag")

	todoAddCmd.Flags().StringVar(&todoFile, "file", "", "Specify the target file")
	todoAddCmd.Flags().StringVar(&todoSchedule, "schedule", "", "Set a SCHEDULED timestamp")
	todoAddCmd.Flags().StringVar(&todoDeadline, "deadline", "", "Set a DEADLINE timestamp")
	todoAddCmd.Flags().StringVar(&todoTags, "tags", "", "Set tags (comma-separated)")
}
