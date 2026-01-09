package config

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestResolveOrgFiles(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "org-agenda-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Define files to create
	filesToCreate := []string{
		"root.org",
		"nested/sub.org",
		"nested/deep/deep.org",
		"other.txt",   // Should be ignored
		".hidden.org", // Should be included if not ignored by logic
	}

	expectedFiles := []string{}
	for _, f := range filesToCreate {
		path := filepath.Join(tmpDir, f)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(path, []byte("test"), 0644)
		if err != nil {
			t.Fatal(err)
		}
		if filepath.Ext(f) == ".org" {
			expectedFiles = append(expectedFiles, path)
		}
	}

	// Run ResolveOrgFiles
	resolved := ResolveOrgFiles([]string{tmpDir})

	// Compare results
	if len(resolved) != 4 { // root.org, sub.org, deep.org, .hidden.org
		t.Errorf("Expected 4 files, got %d", len(resolved))
	}

	sort.Strings(resolved)
	sort.Strings(expectedFiles)

	for i := range resolved {
		if resolved[i] != expectedFiles[i] {
			t.Errorf("Expected %s, got %s", expectedFiles[i], resolved[i])
		}
	}
}

func TestResolveOrgFilesWithMultiplePaths(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "org-agenda-test-multi-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	file1 := filepath.Join(tmpDir, "file1.org")
	file2 := filepath.Join(tmpDir, "file2.org")
	_ = os.WriteFile(file1, []byte("test"), 0644)
	_ = os.WriteFile(file2, []byte("test"), 0644)

	resolved := ResolveOrgFiles([]string{file1, file2})

	if len(resolved) != 2 {
		t.Errorf("Expected 2 files, got %d", len(resolved))
	}
}
