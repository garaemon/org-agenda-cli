package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestAgendaRecursive(t *testing.T) {
	// Setup temporary directory structure
	tmpDir, err := os.MkdirTemp("", "org-agenda-cmd-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	nestedDir := filepath.Join(tmpDir, "nested")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatal(err)
	}

	today := time.Now().Format("2006-01-02")
	dayName := time.Now().Format("Mon")
	timestamp := fmt.Sprintf("<%s %s>", today, dayName)

	rootContent := fmt.Sprintf("* TODO Root Task\nSCHEDULED: %s\n", timestamp)
	nestedContent := fmt.Sprintf("* TODO Nested Task\nSCHEDULED: %s\n", timestamp)

	rootFile := filepath.Join(tmpDir, "root.org")
	nestedFile := filepath.Join(nestedDir, "nested.org")

	if err := os.WriteFile(rootFile, []byte(rootContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nestedFile, []byte(nestedContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Configure viper
	viper.Reset()
	viper.Set("org_files", []string{tmpDir})

	// Reset flags just in case
	agendaDate = ""
	agendaRange = "day"
	agendaTui = false

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	// agendaCmd.Run expects a cobra command and args, but args are not used here
	agendaCmd.Run(agendaCmd, []string{})

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
	if !strings.Contains(output, "Root Task") {
		t.Errorf("Output should contain 'Root Task', got:\n%s", output)
	}
	if !strings.Contains(output, "Nested Task") {
		t.Errorf("Output should contain 'Nested Task', got:\n%s", output)
	}
	if !strings.Contains(output, today) {
		t.Errorf("Output should contain today's date %s, got:\n%s", today, output)
	}
}

func TestAgendaRangeMonth(t *testing.T) {
	// Setup temporary directory structure
	tmpDir, err := os.MkdirTemp("", "org-agenda-cmd-test-month-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a task scheduled for today + 20 days
	future := time.Now().AddDate(0, 0, 20)
	futureDate := future.Format("2006-01-02")
	futureDayName := future.Format("Mon")
	timestamp := fmt.Sprintf("<%s %s>", futureDate, futureDayName)

	// Create a task scheduled for today + 40 days (should not be included)
	farFuture := time.Now().AddDate(0, 0, 40)
	farFutureDate := farFuture.Format("2006-01-02")
	farFutureDayName := farFuture.Format("Mon")
	farTimestamp := fmt.Sprintf("<%s %s>", farFutureDate, farFutureDayName)

	content := fmt.Sprintf("* TODO Near Task\nSCHEDULED: %s\n* TODO Far Task\nSCHEDULED: %s\n", timestamp, farTimestamp)

	file := filepath.Join(tmpDir, "test.org")

	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Configure viper
	viper.Reset()
	viper.Set("org_files", []string{tmpDir})

	// Reset flags
	agendaDate = ""
	agendaRange = "month"
	agendaTui = false

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	agendaCmd.Run(agendaCmd, []string{})

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
	if !strings.Contains(output, "Near Task") {
		t.Errorf("Output should contain 'Near Task', got:\n%s", output)
	}
	if strings.Contains(output, "Far Task") {
		t.Errorf("Output should NOT contain 'Far Task', got:\n%s", output)
	}
}