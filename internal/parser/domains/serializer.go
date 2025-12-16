package domains

import (
	"encoding/json"

	"github.com/unarya/unarya/internal/parser"
)

// SerializeToJSON converts ParseResult into a JSON-serializable structure.
func SerializeToJSON(result *parser.ParseResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SerializeASTAsTree returns a compact tree representation of AST.
func SerializeASTAsTree(node *parser.ASTNode, depth int) string {
	if node == nil {
		return ""
	}
	prefix := ""
	for i := 0; i < depth; i++ {
		prefix += "  "
	}
	out := prefix + node.Type + ": " + node.Name + "\n"
	for _, child := range node.Children {
		out += SerializeASTAsTree(child, depth+1)
	}
	return out
}
