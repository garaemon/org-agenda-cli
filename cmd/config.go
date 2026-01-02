package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manages the configuration file",
	Long:  `Manages the configuration file.`, // Corrected: Removed unnecessary backticks around the string literal
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Display current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		orgFiles := viper.GetStringSlice("org_files")
		fmt.Println("Org Files:")
		for _, f := range orgFiles {
			fmt.Printf("  - %s\n", f)
		}
		fmt.Printf("Default File: %s\n", viper.GetString("default_file"))
		fmt.Printf("Config file used: %s\n", viper.ConfigFileUsed())
	},
}

var configAddPathCmd = &cobra.Command{
	Use:   "add-path [path]",
	Short: "Add an Org file path to the search/display list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Printf("Error resolving path: %v\n", err)
			return
		}

		orgFiles := viper.GetStringSlice("org_files")
		// Check if already exists
		for _, f := range orgFiles {
			if f == path {
				fmt.Printf("Path %s is already in the list.\n", path)
				return
			}
		}

		orgFiles = append(orgFiles, path)
		viper.Set("org_files", orgFiles)

		err = saveConfig()
		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Added path: %s\n", path)
	},
}

var configRemovePathCmd = &cobra.Command{
	Use:   "remove-path [path]",
	Short: "Remove an Org file path from the search/display list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Printf("Error resolving path: %v\n", err)
			return
		}

		orgFiles := viper.GetStringSlice("org_files")
		newOrgFiles := []string{}
		found := false
		for _, f := range orgFiles {
			if f == path {
				found = true
				continue
			}
			newOrgFiles = append(newOrgFiles, f)
		}

		if !found {
			fmt.Printf("Path %s not found in the list.\n", path)
			return
		}

		viper.Set("org_files", newOrgFiles)

		err = saveConfig()
		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Removed path: %s\n", path)
	},
}

func saveConfig() error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configDir := filepath.Join(home, ".config", "org-agenda-cli")
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			err = os.MkdirAll(configDir, 0755)
			if err != nil {
				return err
			}
		}
		configFile = filepath.Join(configDir, "config.yaml")
		return viper.WriteConfigAs(configFile)
	}
	return viper.WriteConfig()
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configAddPathCmd)
	configCmd.AddCommand(configRemovePathCmd)
}
