package security_scan

import (
	"encoding/json"
	"fmt"
)

// GenerateSecurityReport combines all scan results into a unified report
func GenerateSecurityReport(vulns []Vulnerability, secrets []SecretFinding, deps []DependencyIssue, comp ComplianceResult) (*SecurityReport, error) {
	score := calculateRiskScore(vulns, secrets, deps, comp)

	report := &SecurityReport{
		Vulnerabilities: vulns,
		Secrets:         secrets,
		Dependencies:    deps,
		Compliance:      comp,
		RiskScore:       score,
	}

	data, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(data))

	return report, nil
}

func calculateRiskScore(vulns []Vulnerability, secrets []SecretFinding, deps []DependencyIssue, comp ComplianceResult) float64 {
	score := 0.0
	for _, v := range vulns {
		switch v.Severity {
		case "critical":
			score += 3.0
		case "high":
			score += 2.0
		case "medium":
			score += 1.0
		}
	}
	score += float64(len(secrets)) * 2
	score += float64(len(deps))
	score -= float64(len(comp.PassedRules)) * 0.5
	if len(comp.FailedRules) > 0 {
		score += float64(len(comp.FailedRules))
	}
	return score
}
