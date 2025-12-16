package services

import (
	"github.com/unarya/unarya/internal/parser"
	"github.com/unarya/unarya/internal/parser/domains"
)

type ParserService struct{}

func New() *ParserService {
	return &ParserService{}
}

func (p *ParserService) BuildAST(filePath string) (*parser.ASTNode, error) {
	return domains.BuildAST(filePath)
}

func (p *ParserService) ExtractDependencies(rootPath string) ([]parser.Dependency, error) {
	return domains.ExtractDependencies(rootPath)
}

func (p *ParserService) DetectLanguage(filePath string) (lang, version string) {
	return domains.DetectLanguage(filePath)
}

func (p *ParserService) SerializeToJSON(result *parser.ParseResult) (string, error) {
	return domains.SerializeToJSON(result)
}

func (p *ParserService) SerializeASTAsTree(node *parser.ASTNode, depth int) string {
	return domains.SerializeASTAsTree(node, depth)
}

func (p *ParserService) AnalyzeStructure(rootPath string) parser.ProjectStructure {
	return domains.AnalyzeStructure(rootPath)
}
