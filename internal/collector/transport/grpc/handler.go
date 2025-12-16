package grpc

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/unarya/unarya/internal/collector/services"
	pb "github.com/unarya/unarya/lib/proto/pb/collectorpb"
	"github.com/unarya/unarya/pkg/types"
)

type Handler struct {
	pb.UnimplementedCollectorServiceServer
	svc *services.CollectorService
}

func NewCollectorHandler() *Handler {
	return &Handler{
		svc: services.New(),
	}
}

func (h *Handler) CollectFromGit(ctx context.Context, req *pb.GitRequest) (*pb.CollectorResponse, error) {
	// Parse the original URL
	u, err := url.Parse(req.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid git url: %w", err)
	}

	// Only inject token if token existing for private repos
	originalURL := req.Url
	if req.Token != "" {
		// Inject token as basic auth (GitHub/GitLab compatible)
		u.User = url.UserPassword("oauth2", req.Token)
		originalURL = u.String()
	}

	// Determine local path for cloning
	urlForPath, _ := url.Parse(req.Url)
	urlForPath.User = nil // Remove any existing auth
	path := filepath.Join("", urlForPath.Host, strings.TrimSuffix(urlForPath.Path, ".git"))

	gitConfig := types.SourceConfig{
		Ctx:    ctx,
		Type:   "git",
		URL:    originalURL,
		Branch: req.Branch,
		Token:  req.Token,
		Path:   path,
	}

	if err := h.svc.ValidateSource(gitConfig); err != nil {
		return nil, fmt.Errorf("invalid source: %w", err)
	}

	result, err := h.svc.CollectFromGit(gitConfig)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO]: Cloned repository: %s (branch: %s)", req.Url, req.Branch)

	return &pb.CollectorResponse{
		Message: "ok",
		Path:    result.Dir,
	}, nil
}

func (h *Handler) Ready(ctx context.Context, r *pb.Empty) (*pb.ServiceReadyResponse, error) {
	return &pb.ServiceReadyResponse{
		Status: "ok",
	}, nil
}

func (h *Handler) RemoveDir(ctx context.Context, r *pb.Path) (*pb.ServiceOKResponse, error) {
	os.RemoveAll(r.Name)
	return &pb.ServiceOKResponse{
		Message: "ok",
	}, nil
}
