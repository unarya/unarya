package parser

import (
	"path/filepath"
	"strings"
)

// DetectLanguage infers the programming language and version from file extension or content.
func DetectLanguage(filePath string) (lang, version string) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return "Go", "1.x"
	case ".py":
		return "Python", "3.x"
	case ".js":
		return "JavaScript", "ES6+"
	case ".ts":
		return "TypeScript", "4.x"
	case ".java":
		return "Java", "17+"
	case ".cpp", ".cc", ".cxx":
		return "C++", "11+"
	case ".c":
		return "C", "99"
	case ".rs":
		return "Rust", "1.x"
	case ".php":
		return "PHP", "8.x"
	case ".rb":
		return "Ruby", "3.x"
	default:
		return "Unknown", ""
	}
}
