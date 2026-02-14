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

	// Disable TUI for testing
	todoNoInteractive = true

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

func TestTodoAddPriority(t *testing.T) {
	// Setup temporary org file
	tmpfile, err := os.CreateTemp("", "test*.org")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()
	tmpfilePath := tmpfile.Name()
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Reset flags and set for this test
	todoPriority = "A"
	todoFile = tmpfilePath

	// Execute the command logic directly since we set the flags manually
	// We use the function literal from todoAddCmd.Run
	todoAddCmd.Run(todoAddCmd, []string{"Priority Task"})

	// Verify file content
	content, err := os.ReadFile(tmpfilePath)
	if err != nil {
		t.Fatal(err)
	}

	expected := "* TODO [#A] Priority Task"
	if !strings.Contains(string(content), expected) {
		t.Errorf("File content should contain %q, got: %s", expected, string(content))
	}
}

func TestTodoListPriority(t *testing.T) {
	// Setup temporary org file with priorities
	content := `* TODO [#A] Urgent Task
* TODO [#B] Normal Task
* TODO [#C] Low Task
* TODO No Priority Task
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

	// Configure viper
	viper.Reset()
	viper.Set("org_files", []string{tmpfile.Name()})

	// Disable TUI
	todoNoInteractive = true

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	todoListCmd.Run(todoListCmd, []string{})

	// Restore stdout
	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Verify output
	// Expected format: "[TODO] title (file:line)" - currently.
	// We want to see if priority is somehow displayed or at least present.
	// Since we haven't implemented display logic yet, we might see:
	// "[TODO] Urgent Task ..." because Parser strips priority from Title now!
	// Wait, `Parser` extracts Priority. So `Item.Title` does NOT contain `[#A]`.
	// So `todo list` (plain text) currently prints just `Item.Title`.
	// Meaning `[#A]` will be MISSING from output if we don't update `todo list` logic!

	// So this test serves to verify that we NEED to update display logic.
	// Expected: Output should contain "[#A]" or similar if we want it displayed.
	// Current behavior (before fix): It won't contain "[#A]".

	if !strings.Contains(output, "[#A]") {
		t.Logf("Expected output to contain [#A], but it didn't (as expected before implementation). Output: %s", output)
		// We can fail here to drive TDD.
		t.Errorf("Output should contain priority [#A], got: %s", output)
	}
}

func TestTodoListNoColor(t *testing.T) {
	// Setup temporary org file
	content := "* TODO [#A] Urgent Task"
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
	_ = tmpfile.Close()

	viper.Reset()
	viper.Set("org_files", []string{tmpfile.Name()})

	todoNoInteractive = true
	todoNoColor = true // Enable no-color

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run command
	todoListCmd.Run(todoListCmd, []string{})

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Verify output
	// Should contain "[#A]"
	if !strings.Contains(output, "[#A]") {
		t.Errorf("Output should contain priority [#A], got: %s", output)
	}
	// Should NOT contain ANSI codes (simple check for ESC character)
	if strings.Contains(output, "\x1b[") {
		t.Errorf("Output should NOT contain ANSI color codes when no-color is set, got: %q", output)
	}
}

func TestTodoListTags(t *testing.T) {
	// Setup temporary org file with tags
	content := "* TODO Task with Tags :work:urgent:"
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
	_ = tmpfile.Close()

	viper.Reset()
	viper.Set("org_files", []string{tmpfile.Name()})

	todoNoInteractive = true
	todoNoColor = true // simple check for content first

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run command
	todoListCmd.Run(todoListCmd, []string{})

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Verify output
	// Expected to contain ":work:urgent:"
	if !strings.Contains(output, ":work:urgent:") {
		t.Errorf("Output should contain tags :work:urgent:, got: %s", output)
	}
}

func TestTodoListJSON(t *testing.T) {
	// Setup temporary org file
	content := `* TODO [#A] JSON Task :api:`
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
	_ = tmpfile.Close()

	viper.Reset()
	viper.Set("org_files", []string{tmpfile.Name()})

	// Set flags
	todoNoInteractive = true
	todoJSON = true

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run command
	todoListCmd.Run(todoListCmd, []string{})

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Verify output is valid JSON
	// It should be an array of Items
	if !strings.HasPrefix(strings.TrimSpace(output), "[") {
		t.Errorf("Output should start with '[', got: %s", output)
	}

	if !strings.Contains(output, "\"title\": \"JSON Task\"") {
		t.Errorf("Output should contain title field, got: %s", output)
	}
	if !strings.Contains(output, "\"priority\": \"A\"") {
		t.Errorf("Output should contain priority field, got: %s", output)
	}
	if !strings.Contains(output, "\"tags\":") {
		t.Errorf("Output should contain tags field, got: %s", output)
	}
}
