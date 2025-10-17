package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/unarya/unarya/lib/proto/pb/security_scanpb"

	"google.golang.org/grpc"
)

// SecurityScannerServer implements security_scanpb.SecurityScanServiceServer
type SecurityScannerServer struct {
	security_scanpb.UnimplementedSecurityScanServiceServer
}

// main starts the gRPC security scanner service
func main() {
	port := os.Getenv("SECURITY_SCAN_PORT")
	if port == "" {
		port = "50054"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("❌ Failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	security_scanpb.RegisterSecurityScanServiceServer(s, &SecurityScannerServer{})

	log.Printf("🛡️  Security Scan service started on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// ScanForVulnerabilities performs static code analysis
func (s *SecurityScannerServer) ScanForVulnerabilities(ctx context.Context, req *security_scanpb.ScanRequest) (*security_scanpb.ScanResponse, error) {
	sourcePath := req.SourcePath
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("source path not found: %s", sourcePath)
	}

	log.Printf("🔍 Scanning source at %s for vulnerabilities", sourcePath)

	secrets := DetectSecrets(sourcePath)
	depIssues := CheckDependencies(sourcePath)
	permIssues := ValidatePermissions(sourcePath)
	vulnPatterns := DetectCommonVulns(sourcePath)
	report := GenerateSecurityReport(secrets, depIssues, permIssues, vulnPatterns)

	return &security_scanpb.ScanResponse{
		Report:     report,
		TotalFinds: int32(len(secrets) + len(depIssues) + len(permIssues) + len(vulnPatterns)),
	}, nil
}

// DetectSecrets scans files for hardcoded secrets
func DetectSecrets(sourcePath string) []string {
	var secrets []string
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(api[_-]?key|secret|token|password)["'\s:=]+[A-Za-z0-9-_]{8,}`),
		regexp.MustCompile(`(?i)(aws_access_key_id|aws_secret_access_key)\s*=\s*[A-Za-z0-9/+]{20,}`),
	}

	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || strings.Contains(path, "vendor") {
			return nil
		}
		data, _ := os.ReadFile(path)
		for _, p := range patterns {
			if matches := p.FindAllString(string(data), -1); len(matches) > 0 {
				secrets = append(secrets, fmt.Sprintf("%s: %v", path, matches))
			}
		}
		return nil
	})

	return secrets
}

// CheckDependencies scans for known vulnerable dependencies
func CheckDependencies(sourcePath string) []string {
	var issues []string
	depFiles := []string{"go.mod", "package.json", "requirements.txt", "Cargo.toml", "Gemfile.lock"}

	for _, f := range depFiles {
		fp := filepath.Join(sourcePath, f)
		if _, err := os.Stat(fp); err == nil {
			data, _ := os.ReadFile(fp)
			content := string(data)
			if strings.Contains(content, "0.0.1") || strings.Contains(content, "beta") {
				issues = append(issues, fmt.Sprintf("⚠️ Possible vulnerable dependency in %s", f))
			}
		}
	}
	return issues
}

// ValidatePermissions checks file permissions and insecure configs
func ValidatePermissions(sourcePath string) []string {
	var perms []string
	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			mode := info.Mode().Perm()
			if mode&0002 != 0 { // world-writable
				perms = append(perms, fmt.Sprintf("❗ Insecure permission: %s (%#o)", path, mode))
			}
		}
		return nil
	})
	return perms
}

// DetectCommonVulns scans for SQLi, XSS, CSRF patterns
func DetectCommonVulns(sourcePath string) []string {
	var vulns []string
	patterns := map[string]*regexp.Regexp{
		"SQL Injection": regexp.MustCompile(`(?i)SELECT\s+.*\+\s+`),
		"XSS":           regexp.MustCompile(`(?i)innerHTML\s*=`),
		"CSRF":          regexp.MustCompile(`(?i)csrf_token.*missing`),
	}

	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		data, _ := os.ReadFile(path)
		for name, re := range patterns {
			if re.Match(data) {
				vulns = append(vulns, fmt.Sprintf("%s pattern found in %s", name, path))
			}
		}
		return nil
	})
	return vulns
}

// GenerateSecurityReport compiles all issues into a JSON report
func GenerateSecurityReport(secrets, depIssues, perms, vulns []string) string {
	report := map[string]interface{}{
		"secrets":         secrets,
		"dependencies":    depIssues,
		"permissions":     perms,
		"vulnerabilities": vulns,
		"severity": map[string]int{
			"critical": len(secrets),
			"high":     len(vulns),
			"medium":   len(depIssues),
			"low":      len(perms),
		},
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%v"}`, err)
	}
	return string(data)
}
