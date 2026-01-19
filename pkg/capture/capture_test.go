package capture

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInsert_Append(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.org")
	os.WriteFile(file, []byte("* Existing\n"), 0644)

	err := Insert(file, "", nil, "* New\n", false)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	content, _ := os.ReadFile(file)
	expected := "* Existing\n* New\n"
	if string(content) != expected {
		t.Errorf("Expected %q, got %q", expected, string(content))
	}
}

func TestInsert_Prepend(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.org")
	os.WriteFile(file, []byte("* Existing\n"), 0644)

	err := Insert(file, "", nil, "* New\n", true)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	content, _ := os.ReadFile(file)
	expected := "* New\n* Existing\n"
	if string(content) != expected {
		t.Errorf("Expected %q, got %q", expected, string(content))
	}
}

func TestInsert_Heading(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.org")
	content := `* Target
** Child 1
* Other
`
	os.WriteFile(file, []byte(content), 0644)

	// Insert under "Target". Should go after Child 1, before Other.
	// Adjusted level: Target is 1. Child is 2. Entry is "* Entry" (1).
	// Should become "** Entry" (2).

	err := Insert(file, "Target", nil, "* Entry\n", false)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	newContent, _ := os.ReadFile(file)
	expected := `* Target
** Child 1
** Entry
* Other
`
	if string(newContent) != expected {
		t.Errorf("Expected \n%q, got \n%q", expected, newContent)
	}
}

func TestInsert_Heading_Prepend(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.org")
	content := `* Target
** Child 1
`
	os.WriteFile(file, []byte(content), 0644)

	err := Insert(file, "Target", nil, "* Entry\n", true)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	newContent, _ := os.ReadFile(file)
	expected := `* Target
** Entry
** Child 1
`
	if string(newContent) != expected {
		t.Errorf("Expected \n%q, got \n%q", expected, newContent)
	}
}

func TestInsert_OLP(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.org")
	content := `* Level 1
** Level 2
*** Level 3
** Other Level 2
`
	os.WriteFile(file, []byte(content), 0644)

	// Insert under Level 1 -> Level 2
	err := Insert(file, "", []string{"Level 1", "Level 2"}, "* Entry\n", false)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	newContent, _ := os.ReadFile(file)
	expected := `* Level 1
** Level 2
*** Level 3
*** Entry
** Other Level 2
`
	if string(newContent) != expected {
		t.Errorf("Expected \n%q, got \n%q", expected, newContent)
	}
}

func TestInsert_OLP_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.org")
	content := `* Level 1
`
	os.WriteFile(file, []byte(content), 0644)

	err := Insert(file, "", []string{"Level 1", "NonExistent"}, "* Entry\n", false)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInsert_TextOnly(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.org")
	content := `* Target
  Description
** Child 1
`
	os.WriteFile(file, []byte(content), 0644)

	// Insert text under "Target". Should go after Description, before Child 1.
	err := Insert(file, "Target", nil, "New Text\n", false)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	newContent, _ := os.ReadFile(file)
	expected := `* Target
  Description
New Text
** Child 1
`
	if string(newContent) != expected {
		t.Errorf("Expected \n%q, got \n%q", expected, string(newContent))
	}
}

func TestAdjustEntryLevel(t *testing.T) {
	entry := "* Item\n  Body"
	// Target level 2 (inserting under level 1).
	// Entry level 1. Shift +1.
	// Result "** Item\n  Body"

	res := adjustEntryLevel(entry, 2)
	expected := "** Item\n  Body"
	if res != expected {
		t.Errorf("Expected %q, got %q", expected, res)
	}
}

func TestAdjustEntryLevel_Multiline(t *testing.T) {
    entry := "* Item\n** SubItem\n  Body"
    // Shift +1
    res := adjustEntryLevel(entry, 2)
    expected := "** Item\n*** SubItem\n  Body"
    if res != expected {
        t.Errorf("Expected %q, got %q", expected, res)
    }
}
