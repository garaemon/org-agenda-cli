package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/garaemon/org-agenda-cli/pkg/config"
	"github.com/garaemon/org-agenda-cli/pkg/item"
	"github.com/garaemon/org-agenda-cli/pkg/parser"
	"github.com/garaemon/org-agenda-cli/pkg/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	todoStatus        string
	todoTag           string
	todoFile          string
	todoSchedule      string
	todoDeadline      string
	todoTags          string
	todoPriority      string
	todoNoInteractive bool
	todoNoColor       bool
	todoJSON          bool
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
		var allItems []*item.Item
		for _, file := range orgFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file, err)
				continue
			}

			items := parser.ParseString(string(content), file)
			for _, item := range items {
				// If a specific status is requested, filter by it.
				// Otherwise, skip items without any status (i.e., non-TODO headlines)
				// to ensure the 'todo' command only lists actual tasks.
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
				allItems = append(allItems, item)
			}
		}

		if todoJSON {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(allItems); err != nil {
				fmt.Printf("Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
			return
		}

		useTui := !todoNoInteractive

		if useTui {
			if err := tui.Run(allItems, time.Time{}, "", "Todo List"); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}

		// Define styles
		var (
			stylePriorityA = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true) // Red
			stylePriorityB = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))            // Yellow
			stylePriorityC = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))             // Green
			styleStatus    = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true) // Pink
			styleFile      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))            // Grey
		)

		for _, item := range allItems {
			statusStr := item.Status
			if statusStr == "" {
				statusStr = "NONE"
			}

			priorityStr := ""
			if item.Priority != "" {
				priorityStr = fmt.Sprintf("[#%s] ", item.Priority)
			}

			titleStr := item.Title
			fileStr := fmt.Sprintf("(%s:%d)", item.FilePath, item.LineNumber)

			line := ""
			if todoNoColor {
				line = fmt.Sprintf("[%s] %s%s %s", statusStr, priorityStr, titleStr, fileStr)
			} else {
				// Apply colors
				pStyle := lipgloss.NewStyle()
				switch item.Priority {
				case "A":
					pStyle = stylePriorityA
				case "B":
					pStyle = stylePriorityB
				case "C":
					pStyle = stylePriorityC
				}

				sStr := styleStatus.Render("[" + statusStr + "]")
				pStr := pStyle.Render(priorityStr)
				fStr := styleFile.Render(fileStr)

				line = fmt.Sprintf("%s %s%s %s", sStr, pStr, titleStr, fStr)
			}

			// Append tags
			if len(item.Tags) > 0 {
				tagsStr := fmt.Sprintf(" :%s:", strings.Join(item.Tags, ":"))
				if todoNoColor {
					line += tagsStr
				} else {
					// Use standard cyan (6) instead of 86 for better compatibility
					line += lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render(tagsStr)
				}
			}

			fmt.Println(line)
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
		if todoPriority != "" {
			// Validate priority? Org-mode allows any character, but A-C is standard.
			// For now, just insert it.
			// Format: * TODO [#A] Title
			content = fmt.Sprintf("* TODO [#%s] %s", todoPriority, title)
		}

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
		arg := args[0]
		pos, err := parser.ParseFilePosition(arg)
		if err != nil {
			fmt.Printf("Error parsing argument: %v\n", err)
			return
		}

		filePath := pos.FilePath
		lineIdx := pos.Line

		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			return
		}

		lines := strings.Split(string(content), "\n")
		if lineIdx < 1 || lineIdx > len(lines) {
			fmt.Printf("Line number %d is out of range for file %s\n", lineIdx, filePath)
			return
		}

		// Adjust for 0-based index
		targetLine := lines[lineIdx-1]

		// Check if it's a TODO item
		// We expect a line starting with * and containing TODO
		// Regex might be safer, but let's try simple replacement first as per plan "Change TODO to DONE"
		// However, we should be careful not to replace "TODO" in the title.
		// The format is usually "* TODO Title..." or "* TODO [...]"
		// Let's use the parser's regex logic or similar string manipulation.
		// A simple approach: Replace first occurrence of " TODO " with " DONE " after the initial asterisks.

		if !strings.Contains(targetLine, item.StatusTodo) {
			fmt.Printf("Line does not appear to be a %s item.\n", item.StatusTodo)
			return
		}

		// more robust replacement: look for "* TODO" or "* ... TODO"
		// But usually it is "* TODO" or "** TODO".
		// Let's replace " TODO " with " DONE ".
		// If the task status is immediately after stars, it might be "* TODO".

		newLine := strings.Replace(targetLine, " "+item.StatusTodo+" ", " "+item.StatusDone+" ", 1)

		// If no change happened, maybe it's because of strict spacing?
		if newLine == targetLine {
			fmt.Printf("Could not find ' %s ' pattern to replace.\n", item.StatusTodo)
			return
		}

		lines[lineIdx-1] = newLine
		newContent := strings.Join(lines, "\n")

		if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", filePath, err)
			return
		}

		fmt.Printf("Marked task as DONE: %s\n", filePath)
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)

	todoCmd.AddCommand(todoListCmd)
	todoCmd.AddCommand(todoAddCmd)
	todoCmd.AddCommand(todoDoneCmd)

	todoListCmd.Flags().StringVar(&todoStatus, "status", "", "Filter by status (TODO|WAITING|DONE)")
	todoListCmd.Flags().StringVar(&todoTag, "tag", "", "Filter by tag")
	todoListCmd.Flags().BoolVar(&todoNoInteractive, "no-pager", false, "Disable interactive TUI mode")
	todoListCmd.Flags().BoolVar(&todoNoColor, "no-color", false, "Disable colored output")
	todoListCmd.Flags().BoolVar(&todoJSON, "json", false, "Output in JSON format")

	todoAddCmd.Flags().StringVar(&todoFile, "file", "", "Specify the target file")
	todoAddCmd.Flags().StringVar(&todoSchedule, "schedule", "", "Set a SCHEDULED timestamp")
	todoAddCmd.Flags().StringVar(&todoDeadline, "deadline", "", "Set a DEADLINE timestamp")
	todoAddCmd.Flags().StringVar(&todoTags, "tags", "", "Set tags (comma-separated)")
	todoAddCmd.Flags().StringVar(&todoPriority, "priority", "", "Set priority (A, B, C)")
}
