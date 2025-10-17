package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/unarya/unarya/lib/proto/pb/parserpb"
	"google.golang.org/grpc"
)

// ===============================================
// Parser Service — parses source code & extracts structure
// ===============================================

type ParserServer struct {
	parserpb.UnimplementedParserServiceServer
}

// ===============================================
// gRPC Server Entry
// ===============================================
func main() {
	port := os.Getenv("PARSER_PORT")
	if port == "" {
		port = "50053"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("❌ Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	parserpb.RegisterParserServiceServer(grpcServer, &ParserServer{})

	log.Printf("🚀 Parser service started on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}

// ===============================================
// Core RPC Handler
// ===============================================
func (p *ParserServer) ParseCode(ctx context.Context, req *parserpb.ParseRequest) (*parserpb.ParseResponse, error) {
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

	log.Printf("🧩 [Parser] Parsing source directory: %s", sourcePath)

	// 1️⃣ Detect language
	lang := DetectLanguage(sourcePath)
	log.Printf("🗣️  Detected language: %s", lang)

	// 2️⃣ Extract dependencies
	deps := ExtractDependencies(sourcePath)
	log.Printf("📦 Found %d dependency files", len(deps))

	// 3️⃣ Build AST / code structure
	astData := BuildAST(sourcePath, lang)

	// 4️⃣ Convert to JSON structure
	codeStructure := GenerateCodeRepresentation(astData)

	resp := &parserpb.ParseResponse{
		Language:       lang,
		Dependencies:   deps,
		CodeStructure:  codeStructure,
		Representation: "json",
	}

	log.Printf("✅ [Parser] Completed parsing (%s)", lang)
	return resp, nil
}

// ===============================================
// Helpers
// ===============================================

// DetectLanguage identifies the main programming language in a directory
func DetectLanguage(sourcePath string) string {
	files, err := os.ReadDir(sourcePath)
	if err != nil {
		return "Unknown"
	}

	extCount := make(map[string]int)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name())
		extCount[ext]++
	}

	// Pick the most common extension
	mainExt := ""
	maxCount := 0
	for ext, count := range extCount {
		if count > maxCount {
			maxCount = count
			mainExt = ext
		}
	}

	switch mainExt {
	case ".go":
		return "Go"
	case ".py":
		return "Python"
	case ".js":
		return "JavaScript"
	case ".ts":
		return "TypeScript"
	case ".java":
		return "Java"
	case ".cpp", ".c", ".hpp":
		return "C/C++"
	case ".rs":
		return "Rust"
	case ".php":
		return "PHP"
	case ".rb":
		return "Ruby"
	default:
		return "Unknown"
	}
}

// ExtractDependencies scans common dependency/config files (go.mod, package.json, etc.)
func ExtractDependencies(sourcePath string) []string {
	var deps []string
	depFiles := []string{
		"go.mod", "package.json", "requirements.txt",
		"pom.xml", "Cargo.toml", "composer.json", "Gemfile",
	}

	for _, file := range depFiles {
		fullPath := filepath.Join(sourcePath, file)
		if data, err := os.ReadFile(fullPath); err == nil {
			deps = append(deps, fmt.Sprintf("%s: %d bytes", file, len(data)))
		}
	}
	return deps
}

// BuildAST builds a basic AST (only implemented for Go for now)
func BuildAST(sourcePath, lang string) interface{} {
	if lang != "Go" {
		return map[string]any{"note": fmt.Sprintf("AST for %s not implemented", lang)}
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, sourcePath, nil, parser.ParseComments)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}

	return pkgs
}

// GenerateCodeRepresentation converts AST or structure into JSON string
func GenerateCodeRepresentation(astData interface{}) string {
	jsonData, err := json.MarshalIndent(astData, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%v"}`, err)
	}
	return string(jsonData)
}
