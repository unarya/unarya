// file: /internal/collector/scan.go
package domain

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/unarya/unarya/pkg/types"
)

// scanFiles recursively walks a directory and gathers file info.
func scanFiles(root string) (*types.CollectionResult, error) {
	result := &types.CollectionResult{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		result.Files = append(result.Files, types.FileInfo{
			Name: filepath.Base(path),
			Path: path,
			Size: info.Size(),
			// language detection could be added later
		})
		result.TotalSize += info.Size()
		return nil
	})
	return result, err
}
