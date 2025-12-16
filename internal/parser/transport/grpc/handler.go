package grpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/unarya/unarya/internal/parser/helpers"
	"github.com/unarya/unarya/internal/parser/services"
	"github.com/unarya/unarya/lib/proto/pb/parserpb"
)

type Handler struct {
	parserpb.UnimplementedParserServiceServer
	svc *services.ParserService
}

func NewParserHandler() *Handler {
	return &Handler{
		svc: services.New(),
	}
}

// ParseCode
//
// Core RPC Handler
// ===============================================
func (p *Handler) ParseCode(ctx context.Context, req *parserpb.ParseRequest) (*parserpb.ParseResponse, error) {
	sourcePath := strings.TrimSpace(req.SourcePath)
	if sourcePath == "" {
		return nil, fmt.Errorf("source_path is empty")
	}

	info, err := os.Stat(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("cannot access source path: %v", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("source path must be a directory: %s", sourcePath)
	}

	log.Printf("[INFO] Parsing source directory: %s", sourcePath)

	// 1. Detect language
	lang := helpers.DetectLanguage(sourcePath)
	log.Printf("[INFO] Detected language: %s", lang)

	// 2. Extract dependencies
	deps := helpers.ExtractDependencies(sourcePath)
	log.Printf("[INFO] Found %d dependency files", len(deps))

	// 3. Build AST / code structure
	astData := helpers.BuildAST(sourcePath, lang)

	// 4. Convert to JSON structure
	codeStructure := helpers.GenerateCodeRepresentation(astData)

	resp := &parserpb.ParseResponse{
		Language:       lang,
		Dependencies:   deps,
		CodeStructure:  codeStructure,
		Representation: "json",
	}

	log.Printf("[Parser] Completed parsing (%s)", lang)
	return resp, nil
}
