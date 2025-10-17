package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/unarya/unarya/lib/proto/pb/collectorpb"
	"google.golang.org/grpc"
)

// CollectorServer implements collectorpb.CollectorServiceServer
type CollectorServer struct {
	collectorpb.UnimplementedCollectorServiceServer
}

// main starts the gRPC Collector service
func main() {
	port := os.Getenv("COLLECTOR_PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	collectorpb.RegisterCollectorServiceServer(s, &CollectorServer{})

	log.Printf("🚀 Collector service started on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// CollectFromGit clones a repository from Git with optional authentication
func (c *CollectorServer) CollectFromGit(ctx context.Context, req *collectorpb.GitRequest) (*collectorpb.CollectorResponse, error) {
	if err := ValidateSource(req.Url); err != nil {
		return nil, fmt.Errorf("invalid source: %w", err)
	}

	dir := filepath.Join(os.TempDir(), fmt.Sprintf("repo-%d", time.Now().UnixNano()))
	cloneCmd := []string{"git", "clone"}

	if req.Branch != "" {
		cloneCmd = append(cloneCmd, "-b", req.Branch)
	}
	if req.Token != "" {
		urlParts := strings.SplitN(req.Url, "://", 2)
		if len(urlParts) == 2 {
			req.Url = fmt.Sprintf("https://%s@%s", req.Token, urlParts[1])
		}
	}

	cloneCmd = append(cloneCmd, req.Url, dir)

	cmd := exec.CommandContext(ctx, cloneCmd[0], cloneCmd[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git clone failed: %v\n%s", err, string(output))
	}

	log.Printf("✅ Cloned repository: %s (branch: %s)", req.Url, req.Branch)
	return &collectorpb.CollectorResponse{
		Message: fmt.Sprintf("Repository cloned successfully at %s", dir),
		Path:    dir,
	}, nil
}

// CollectFromArchive downloads and extracts a ZIP/TAR archive
func (c *CollectorServer) CollectFromArchive(ctx context.Context, req *collectorpb.ArchiveRequest) (*collectorpb.CollectorResponse, error) {
	if err := ValidateSource(req.Url); err != nil {
		return nil, err
	}

	targetDir := filepath.Join(os.TempDir(), fmt.Sprintf("archive-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, err
	}

	// Use system commands for simplicity (curl + tar or unzip)
	log.Printf("📦 Downloading archive from %s", req.Url)
	var cmd *exec.Cmd
	if strings.HasSuffix(req.Url, ".zip") {
		cmd = exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("curl -L %s -o /tmp/archive.zip && unzip /tmp/archive.zip -d %s", req.Url, targetDir))
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("curl -L %s | tar -xz -C %s", req.Url, targetDir))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("archive extraction failed: %v\n%s", err, string(output))
	}

	log.Printf("✅ Archive extracted to %s", targetDir)
	return &collectorpb.CollectorResponse{
		Message: "Archive downloaded and extracted successfully",
		Path:    targetDir,
	}, nil
}

// CollectFromURL downloads raw files from direct URLs
func (c *CollectorServer) CollectFromURL(ctx context.Context, req *collectorpb.URLRequest) (*collectorpb.CollectorResponse, error) {
	if err := ValidateSource(req.Url); err != nil {
		return nil, err
	}

	target := filepath.Join(os.TempDir(), fmt.Sprintf("file-%d", time.Now().UnixNano()))
	cmd := exec.CommandContext(ctx, "curl", "-L", req.Url, "-o", target)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("download failed: %v\n%s", err, string(output))
	}

	log.Printf("✅ File downloaded to %s", target)
	return &collectorpb.CollectorResponse{
		Message: "File downloaded successfully",
		Path:    target,
	}, nil
}

// ValidateSource performs security checks to prevent unsafe URLs or paths
func ValidateSource(url string) error {
	if url == "" {
		return fmt.Errorf("empty URL")
	}

	if strings.Contains(url, "..") {
		return fmt.Errorf("path traversal detected")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "git@") && !strings.HasPrefix(url, "ssh://") {
		return fmt.Errorf("invalid URL scheme")
	}

	return nil
}
