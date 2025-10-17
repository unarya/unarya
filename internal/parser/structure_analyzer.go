package parser

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// AnalyzeStructure inspects the directory and infers architecture patterns.
func AnalyzeStructure(rootPath string) ProjectStructure {
	structure := ProjectStructure{}
	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Detect Go/Node/Python project conventions
			if strings.Contains(strings.ToLower(d.Name()), "pkg") {
				structure.Packages = append(structure.Packages, d.Name())
			}
			if strings.Contains(strings.ToLower(d.Name()), "module") {
				structure.Modules = append(structure.Modules, d.Name())
			}
			return nil
		}

		structure.Files = append(structure.Files, path)

		if strings.HasSuffix(d.Name(), "main.go") || strings.HasSuffix(d.Name(), "index.js") {
			structure.EntryFile = path
		}
		return nil
	})
	return structure
}
