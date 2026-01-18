package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/garaemon/org-agenda-cli/pkg/capture"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var captureFile string

var captureCmd = &cobra.Command{
	Use:   "capture [content]",
	Short: "Capture a note to an Org file",
	Long:  `Capture a note to an Org file using a configurable format.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content := strings.Join(args, " ")

		// Determine target file
		targetFile := captureFile
		if targetFile == "" {
			targetFile = viper.GetString("capture.default_file")
		}
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

		// Apply date formatting to targetFile (e.g. %Y-capture.org -> 2026-capture.org)
		targetFile = capture.Format(targetFile, "")

		// Determine format
		format := viper.GetString("capture.format")
		if format == "" {
			format = "* %t\n  %c"
		}

		// Format entry
		entry := capture.Format(format, content)

		// Ensure newline at the end
		if !strings.HasSuffix(entry, "\n") {
			entry += "\n"
		}

		prepend := viper.GetBool("capture.prepend")

		if prepend {
			// Read existing content
			existingContent, err := os.ReadFile(targetFile)
			if err != nil && !os.IsNotExist(err) {
				fmt.Printf("Error reading file: %v\n", err)
				return
			}

			// Prepend new entry
			newContent := entry
			if len(existingContent) > 0 {
				newContent += string(existingContent)
			}

			// Write back
			if err := os.WriteFile(targetFile, []byte(newContent), 0644); err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		} else {
			// Append to file
			f, err := os.OpenFile(targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Error opening file: %v\n", err)
				return
			}
			defer f.Close()

			if _, err := f.WriteString(entry); err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		}

		fmt.Printf("Captured to %s\n", targetFile)
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)
	captureCmd.Flags().StringVar(&captureFile, "file", "", "Specify the target file")
}
