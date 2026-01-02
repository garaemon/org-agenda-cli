package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigRemovePath(t *testing.T) {
	// Setup temporary config file
	tmpConfigDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpConfigDir)

	configPath := filepath.Join(tmpConfigDir, "config.yaml")
	
	// Create dummy org files to test with
	file1, _ := filepath.Abs("test1.org")
	file2, _ := filepath.Abs("test2.org")

	// Set initial config
	viper.Reset()
	viper.SetConfigFile(configPath)
	viper.Set("org_files", []string{file1, file2})
	if err := viper.WriteConfigAs(configPath); err != nil {
		t.Fatal(err)
	}

	t.Run("Remove existing path", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Run the command to remove file1
		configRemovePathCmd.Run(configRemovePathCmd, []string{"test1.org"})

		// Restore stdout
		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		if !strings.Contains(output, "Removed path") {
			t.Errorf("Expected 'Removed path' in output, got: %s", output)
		}

		// Verify viper's state
		orgFiles := viper.GetStringSlice("org_files")
		if len(orgFiles) != 1 || orgFiles[0] != file2 {
			t.Errorf("Expected org_files to have %s, got %v", file2, orgFiles)
		}

		// Verify file state
		if err := viper.ReadInConfig(); err != nil {
			t.Fatal(err)
		}
		orgFiles = viper.GetStringSlice("org_files")
		if len(orgFiles) != 1 || orgFiles[0] != file2 {
			t.Errorf("Expected org_files in file to have %s, got %v", file2, orgFiles)
		}
	})

	t.Run("Remove non-existing path", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Run the command to remove non-existing file
		configRemovePathCmd.Run(configRemovePathCmd, []string{"nonexistent.org"})

		// Restore stdout
		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		if !strings.Contains(output, "not found in the list") {
			t.Errorf("Expected 'not found in the list' in output, got: %s", output)
		}

		// Verify viper's state remains unchanged
		orgFiles := viper.GetStringSlice("org_files")
		if len(orgFiles) != 1 || orgFiles[0] != file2 {
			t.Errorf("Expected org_files to remain %v, got %v", []string{file2}, orgFiles)
		}
	})
}

func TestConfigAddPath(t *testing.T) {
	// Setup temporary config file
	tmpConfigDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpConfigDir)

	configPath := filepath.Join(tmpConfigDir, "config.yaml")

	// Set initial config
	viper.Reset()
	viper.SetConfigFile(configPath)
	if err := viper.WriteConfigAs(configPath); err != nil {
		t.Fatal(err)
	}

	t.Run("Add new path", func(t *testing.T) {
		file1, _ := filepath.Abs("test1.org")

		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Run the command to add file1
		configAddPathCmd.Run(configAddPathCmd, []string{"test1.org"})

		// Restore stdout
		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		if !strings.Contains(output, "Added path") {
			t.Errorf("Expected 'Added path' in output, got: %s", output)
		}

		// Verify viper's state
		orgFiles := viper.GetStringSlice("org_files")
		if len(orgFiles) != 1 || orgFiles[0] != file1 {
			t.Errorf("Expected org_files to have %s, got %v", file1, orgFiles)
		}
	})

	t.Run("Add existing path", func(t *testing.T) {
		file1, _ := filepath.Abs("test1.org")

		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Run the command to add file1 again
		configAddPathCmd.Run(configAddPathCmd, []string{"test1.org"})

		// Restore stdout
		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		if !strings.Contains(output, "is already in the list") {
			t.Errorf("Expected 'is already in the list' in output, got: %s", output)
		}

		// Verify viper's state
		orgFiles := viper.GetStringSlice("org_files")
		if len(orgFiles) != 1 || orgFiles[0] != file1 {
			t.Errorf("Expected org_files to still have only %s, got %v", file1, orgFiles)
		}
	})
}
