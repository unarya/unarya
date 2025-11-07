package types

import "errors"

// SourceConfig defines the configuration for a source collection request.
type SourceConfig struct {
	Type      string // "git", "archive", "url"
	URL       string
	Branch    string
	Token     string
	LocalPath string
}

// FileInfo represents a collected file's metadata.
type FileInfo struct {
	Name     string
	Path     string
	Size     int64
	Language string
}

// CollectionResult summarizes the result of a collection operation.
type CollectionResult struct {
	Files     []FileInfo
	TotalSize int64
	Language  []string
	Error     error
}

// ErrInvalidSourceType is returned when the source type is unsupported.
var ErrInvalidSourceType = errors.New("invalid source type")

// SupportedTypes returns all acceptable source types.
func SupportedTypes() []string {
	return []string{"git", "archive", "url"}
}
