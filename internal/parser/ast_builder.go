package parser

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// BuildAST constructs an Abstract Syntax Tree for supported languages.
// Currently implemented for Go, can be extended to others.
func BuildAST(filePath string) (*ASTNode, error) {
	lang, _ := DetectLanguage(filePath)
	switch lang {
	case "Go":
		return buildGoAST(filePath)
	default:
		// Placeholder: support future languages
		return &ASTNode{
			Type: "File",
			Name: filepath.Base(filePath),
			Metadata: map[string]interface{}{
				"note": "AST not implemented for this language yet",
			},
		}, nil
	}
}

func buildGoAST(filePath string) (*ASTNode, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	root := &ASTNode{
		Type:     "GoFile",
		Name:     node.Name.Name,
		Metadata: map[string]interface{}{"imports": len(node.Imports)},
	}

	for _, decl := range node.Decls {
		child := &ASTNode{
			Type: "Declaration",
			Metadata: map[string]interface{}{
				"type": fmt.Sprintf("%T", decl),
			},
		}
		root.Children = append(root.Children, child)
	}
	return root, nil
}
