package cmd

import (
	"fmt"
	"os"

	"github.com/garaemon/org-agenda-cli/pkg/mcp"
	"github.com/garaemon/org-agenda-cli/pkg/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
	Long:  `Start the MCP server to expose org-agenda-cli functionality via Model Context Protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		paths := viper.GetStringSlice("org_files")
		defaultFile := viper.GetString("default_file")

		if len(paths) == 0 {
			// Write warnings to stderr to avoid interfering with MCP stdio transport.
			fmt.Fprintln(os.Stderr, "Warning: No org files configured.")
		}

		svc := service.NewService(paths, defaultFile)
		s := mcp.NewServer(svc)

		if err := s.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
