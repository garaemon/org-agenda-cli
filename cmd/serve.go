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
			// Fallback or warning?
			// For MCP server, maybe we want to be silent on stdout if possible,
			// but we are using stdio for transport.
			// However, logs usually go to stderr.
			// Let's print to stderr.
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
