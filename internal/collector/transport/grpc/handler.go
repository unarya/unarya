package grpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/unarya/unarya/internal/collector/service"
	pb "github.com/unarya/unarya/lib/proto/pb/collectorpb"
	"github.com/unarya/unarya/pkg/types"
)

type Handler struct {
	pb.UnimplementedCollectorServiceServer
	svc *service.CollectorService
}

func NewCollectorHandler() *Handler {
	return &Handler{
		svc: service.New(),
	}
}

func (h *Handler) CollectFromGit(ctx context.Context, req *pb.GitRequest) (*pb.CollectorResponse, error) {
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

	gitConfig := types.SourceConfig{
		Type:      "git",
		URL:       req.Url,
		Branch:    req.Branch,
		Token:     req.Token,
		LocalPath: "",
	}
	if err := h.svc.ValidateSource(gitConfig); err != nil {
		return nil, fmt.Errorf("invalid source: %w", err)
	}
	result, err := h.svc.CollectFromGit(types.SourceConfig{
		URL:    req.Url,
		Token:  req.Token,
		Branch: req.Branch,
	})
	if err != nil {
		return nil, err
	}
	// Preprocessing for metadata after cloned git source
	fmt.Printf("Git resource data: %v", result)
	cloneCmd = append(cloneCmd, req.Url, dir)

	cmd := exec.CommandContext(ctx, cloneCmd[0], cloneCmd[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git clone failed: %v\n%s", err, string(output))
	}

	log.Printf("✅ Cloned repository: %s (branch: %s)", req.Url, req.Branch)

	return &pb.CollectorResponse{
		Message: "ok",
		Path:    dir,
	}, nil
}
