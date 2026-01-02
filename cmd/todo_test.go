package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestTodoList(t *testing.T) {
	// Setup temporary org file
	content := `* TODO Task 1
* DONE Task 2
* Non-TODO Item
`
	tmpfile, err := os.CreateTemp("", "test*.org")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Configure viper to use the temp file
	viper.Reset()
	viper.Set("org_files", []string{tmpfile.Name()})

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	todoListCmd.Run(todoListCmd, []string{})

	// Restore stdout
	if err := w.Close(); err != nil {
		t.Logf("failed to close pipe: %v", err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	// Verify output
	if !strings.Contains(output, "[TODO] Task 1") {
		t.Errorf("Output should contain Task 1, got: %s", output)
	}
	if !strings.Contains(output, "[DONE] Task 2") {
		t.Errorf("Output should contain Task 2, got: %s", output)
	}
	if strings.Contains(output, "Non-TODO Item") {
		t.Errorf("Output should NOT contain Non-TODO Item, got: %s", output)
	}
}
