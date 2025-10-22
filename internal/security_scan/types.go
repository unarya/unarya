package security_scan

type SecurityReport struct {
	Vulnerabilities []Vulnerability
	Secrets         []SecretFinding
	Dependencies    []DependencyIssue
	Compliance      ComplianceResult
	RiskScore       float64
}

type Vulnerability struct {
	Type        string
	Severity    string // "critical", "high", "medium", "low"
	Location    Location
	Description string
	Remediation string
}

type SecretFinding struct {
	Type        string
	Severity    string
	Line        int
	FilePath    string
	Pattern     string
	Entropy     float64
	Description string
}

type DependencyIssue struct {
	PackageName string
	Version     string
	CVE         string
	Severity    string
	Description string
	FixedIn     string
}

type ComplianceResult struct {
	PassedRules   []string
	FailedRules   []ComplianceRule
	Warnings      []string
	PolicyVersion string
}

type ComplianceRule struct {
	ID             string
	Description    string
	Severity       string
	Recommendation string
}

type Location struct {
	File string
	Line int
	Col  int
}
