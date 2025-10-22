package security_scan

import (
	"path/filepath"
	"strings"
)

// ScanForVulnerabilities performs static code analysis for known vulnerability patterns
func ScanForVulnerabilities(codeDir string) ([]Vulnerability, error) {
	// Placeholder logic — normally you'd parse AST or regex scan code
	vulns := []Vulnerability{}

	if strings.Contains(codeDir, "sql") {
		vulns = append(vulns, Vulnerability{
			Type:        "SQL Injection",
			Severity:    "high",
			Location:    Location{File: filepath.Join(codeDir, "example.go"), Line: 23},
			Description: "Potential SQL Injection found in query construction",
			Remediation: "Use prepared statements or parameterized queries",
		})
	}

	return vulns, nil
}
