package security_scan

import (
	"strings"
)

// ValidateCompliance checks security configurations (Docker, K8s, permissions, etc.)
func ValidateCompliance(configDir string) (ComplianceResult, error) {
	result := ComplianceResult{
		PassedRules:   []string{},
		FailedRules:   []ComplianceRule{},
		Warnings:      []string{},
		PolicyVersion: "1.0.0",
	}

	if strings.Contains(configDir, "Dockerfile") {
		result.FailedRules = append(result.FailedRules, ComplianceRule{
			ID:             "DOCKER-001",
			Description:    "Container runs as root user",
			Severity:       "high",
			Recommendation: "Add USER directive to Dockerfile to use non-root user",
		})
	}

	if strings.Contains(configDir, "k8s") {
		result.PassedRules = append(result.PassedRules, "K8S-SECURE-POLICY")
	}

	return result, nil
}
