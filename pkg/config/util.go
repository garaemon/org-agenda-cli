package config

import (
	"os"
	"path/filepath"
)

// ResolveOrgFiles takes a list of paths (files or directories) and returns a list of all .org files found.
func ResolveOrgFiles(paths []string) []string {
	var files []string
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.IsDir() {
			_ = filepath.Walk(path, func(p string, i os.FileInfo, e error) error {
				if e == nil && !i.IsDir() && filepath.Ext(p) == ".org" {
					files = append(files, p)
				}
				return nil
			})
		} else {
			if filepath.Ext(path) == ".org" {
				files = append(files, path)
			}
		}
	}
	return files
}
