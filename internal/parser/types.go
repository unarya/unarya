package parser

// ParseResult represents the full outcome of a source code parsing operation.
type ParseResult struct {
	Language     string
	Version      string
	AST          *ASTNode
	Dependencies []Dependency
	Structure    ProjectStructure
	Metrics      CodeMetrics
}

// ASTNode defines a node in the Abstract Syntax Tree.
type ASTNode struct {
	Type     string
	Name     string
	Children []*ASTNode
	Metadata map[string]interface{}
}

// Dependency represents a project dependency (from go.mod, package.json, etc.).
type Dependency struct {
	Name    string
	Version string
	Source  string
}

// ProjectStructure summarizes the organization of files and directories.
type ProjectStructure struct {
	Modules   []string
	Packages  []string
	EntryFile string
	Files     []string
}

// CodeMetrics captures basic statistics about the codebase.
type CodeMetrics struct {
	TotalFiles int
	TotalLines int
	Functions  int
	Classes    int
	Complexity int
}
