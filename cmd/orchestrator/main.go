package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/unarya/unarya/lib/proto/pb/aipb"
	"github.com/unarya/unarya/lib/proto/pb/collectorpb"
	"github.com/unarya/unarya/lib/proto/pb/orchestratorpb"
	"github.com/unarya/unarya/lib/proto/pb/parserpb"
	"github.com/unarya/unarya/lib/proto/pb/security_scanpb"
	"google.golang.org/grpc"
)

// ===============================================
// Orchestrator — central coordinator
// collector → parser → ai → security_scan
// ===============================================

type OrchestratorServer struct {
	orchestratorpb.UnimplementedOrchestratorServiceServer

	collectorClient collectorpb.CollectorServiceClient
	parserClient    parserpb.ParserServiceClient
	aiClient        aipb.AIServiceClient
	securityClient  security_scanpb.SecurityScanServiceClient
}

// StartPipeline — coordinates the full pipeline flow
func (s *OrchestratorServer) StartPipeline(ctx context.Context, req *orchestratorpb.PipelineRequest) (*orchestratorpb.PipelineResponse, error) {
	log.Printf("[Orchestrator] Received pipeline request for repo: %s", req.RepositoryUrl)

	// === 1️⃣ Collector stage ===
	collected, err := s.collectorClient.CollectFromGit(ctx, &collectorpb.GitRequest{
		Url: req.RepositoryUrl,
	})
	if err != nil {
		return s.fail("collector", err)
	}
	log.Println("[Orchestrator] ✓ Repository collected")

	// === 2️⃣ Parser stage ===
	parsed, err := s.parserClient.ParseCode(ctx, &parserpb.ParseRequest{
		SourcePath: collected.Path,
	})
	if err != nil {
		return s.fail("parser", err)
	}
	log.Println("[Orchestrator] ✓ Parsing completed")

	// === 3️⃣ AI analysis stage ===
	analyzed, err := s.aiClient.AnalyzeCode(ctx, &aipb.AIAnalyzeRequest{
		Language:      parsed.Language,
		CodeStructure: parsed.CodeStructure,
	})
	if err != nil {
		return s.fail("ai", err)
	}
	log.Println("[Orchestrator] ✓ AI analysis completed")

	// === 4️⃣ Security scanning stage ===
	scanned, err := s.securityClient.ScanForVulnerabilities(ctx, &security_scanpb.ScanRequest{
		SourcePath: collected.Path,
	})
	if err != nil {
		return s.fail("security_scan", err)
	}
	log.Println("[Orchestrator] ✓ Security scan completed")

	// === 5️⃣ Aggregate results ===
	details := fmt.Sprintf(
		"AI insights: %s (confidence: %s)\nSecurity findings: %d issues\nReport summary: %s",
		analyzed.Insights, analyzed.Confidence, scanned.TotalFinds, scanned.Report,
	)

	return &orchestratorpb.PipelineResponse{
		Status:  "success",
		Details: details,
	}, nil
}

func (s *OrchestratorServer) fail(stage string, err error) (*orchestratorpb.PipelineResponse, error) {
	log.Printf("[ERROR] Stage '%s' failed: %v", stage, err)
	return &orchestratorpb.PipelineResponse{
		Status:  "failed",
		Details: fmt.Sprintf("Stage '%s' failed: %v", stage, err),
	}, err
}

// StartOrchestrator launches the orchestrator gRPC server
func StartOrchestrator() error {
	port := os.Getenv("ORCHESTRATOR_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	orchestratorpb.RegisterOrchestratorServiceServer(grpcServer, &OrchestratorServer{})

	log.Printf("[Unarya] 🚀 Orchestrator service started on port %s", port)
	return grpcServer.Serve(lis)
}

func main() {
	if err := StartOrchestrator(); err != nil {
		log.Fatalf("❌ Failed to start orchestrator: %v", err)
	}
}
