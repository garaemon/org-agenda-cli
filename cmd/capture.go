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

		fmt.Printf("Captured to %s\n", targetFile)
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)
	captureCmd.Flags().StringVar(&captureFile, "file", "", "Specify the target file")
}
