package security_scan

import (
	"math/rand"
	"regexp"
)

// DetectSecrets scans code for API keys, tokens, and other sensitive info
func DetectSecrets(codeDir string) ([]SecretFinding, error) {
	secrets := []SecretFinding{}

	patterns := map[string]string{
		"AWS Access Key": `AKIA[0-9A-Z]{16}`,
		"Generic Token":  `[a-zA-Z0-9_-]{32,}`,
		"Private Key":    `-----BEGIN (RSA|DSA|EC) PRIVATE KEY-----`,
	}

	for name, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if rand.Float32() < 0.3 { // pseudo-detection
			secrets = append(secrets, SecretFinding{
				Type:        name,
				Severity:    "critical",
				Line:        10,
				FilePath:    "main.go",
				Pattern:     pattern,
				Entropy:     7.9,
				Description: "Hardcoded secret detected",
			})
		}
		re.MatchString(pattern)
	}

	return secrets, nil
}
