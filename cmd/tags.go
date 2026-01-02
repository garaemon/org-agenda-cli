package cmd

import (
	"fmt"
	"os"

	"github.com/garaemon/org-agenda-cli/pkg/agenda"
	"github.com/garaemon/org-agenda-cli/pkg/config"
	"github.com/garaemon/org-agenda-cli/pkg/item"
	"github.com/garaemon/org-agenda-cli/pkg/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Lists all unique tags",
	Long:  `Lists all unique tags across all configured Org files.`,
	Run: func(cmd *cobra.Command, args []string) {
		paths := viper.GetStringSlice("org_files")
		if len(paths) == 0 {
			if _, err := os.Stat("sample.org"); err == nil {
				paths = []string{"sample.org"}
			} else {
				fmt.Println("No org files configured.")
				return
			}
		}

		orgFiles := config.ResolveOrgFiles(paths)
		var allItems []*item.Item

		for _, file := range orgFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			items := parser.ParseString(string(content), file)
			allItems = append(allItems, items...)
		}

		tags := agenda.ExtractUniqueTags(allItems)

		if len(tags) == 0 {
			fmt.Println("No tags found.")
			return
		}

		for _, tag := range tags {
			fmt.Println(tag)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagsCmd)
}
