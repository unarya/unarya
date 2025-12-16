package helpers

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// DetectLanguage identifies the main programming language in a directory
func DetectLanguage(sourcePath string) string {
	files, err := os.ReadDir(sourcePath)
	if err != nil {
		return "Unknown"
	}

	extCount := make(map[string]int)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name())
		extCount[ext]++
	}

	// Pick the most common extension
	mainExt := ""
	maxCount := 0
	for ext, count := range extCount {
		if count > maxCount {
			maxCount = count
			mainExt = ext
		}
	}

	switch mainExt {
	case ".go":
		return "Go"
	case ".py":
		return "Python"
	case ".js":
		return "JavaScript"
	case ".ts":
		return "TypeScript"
	case ".java":
		return "Java"
	case ".cpp", ".c", ".hpp":
		return "C/C++"
	case ".rs":
		return "Rust"
	case ".php":
		return "PHP"
	case ".rb":
		return "Ruby"
	default:
		return "Unknown"
	}
}

// ExtractDependencies scans common dependency/config files (go.mod, package.json, etc.)
func ExtractDependencies(sourcePath string) []string {
	var deps []string
	depFiles := []string{
		"go.mod", "package.json", "requirements.txt",
		"pom.xml", "Cargo.toml", "composer.json", "Gemfile",
	}

	for _, file := range depFiles {
		fullPath := filepath.Join(sourcePath, file)
		if data, err := os.ReadFile(fullPath); err == nil {
			deps = append(deps, fmt.Sprintf("%s: %d bytes", file, len(data)))
		}
	}
	return deps
}

// BuildAST builds a basic AST (only implemented for Go for now)
func BuildAST(sourcePath, lang string) interface{} {
	if lang != "Go" {
		return map[string]any{"note": fmt.Sprintf("AST for %s not implemented", lang)}
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, sourcePath, nil, parser.ParseComments)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}

	return pkgs
}

// GenerateCodeRepresentation converts AST or structure into JSON string
func GenerateCodeRepresentation(astData interface{}) string {
	jsonData, err := json.MarshalIndent(astData, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%v"}`, err)
	}
	return string(jsonData)
}
