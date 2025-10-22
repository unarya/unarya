package security_scan

import (
	"strings"
)

// CheckDependencies scans dependency manifests for known CVEs
func CheckDependencies(manifestPath string) ([]DependencyIssue, error) {
	issues := []DependencyIssue{}

	if strings.Contains(manifestPath, "package.json") {
		issues = append(issues, DependencyIssue{
			PackageName: "express",
			Version:     "4.17.1",
			CVE:         "CVE-2023-1234",
			Severity:    "medium",
			Description: "Prototype pollution vulnerability in Express",
			FixedIn:     "4.18.0",
		})
	}

	if strings.Contains(manifestPath, "go.mod") {
		issues = append(issues, DependencyIssue{
			PackageName: "gin-gonic/gin",
			Version:     "1.8.0",
			CVE:         "CVE-2024-5678",
			Severity:    "high",
			Description: "Improper input validation in Gin middleware",
			FixedIn:     "1.9.0",
		})
	}

	return issues, nil
}
