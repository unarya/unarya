package collector

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// CollectFromURL downloads a single file or directory structure via HTTP(S).
func CollectFromURL(cfg SourceConfig) (*CollectionResult, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("url is empty")
	}
	if cfg.LocalPath == "" {
		cfg.LocalPath = filepath.Join(os.TempDir(), "collector_http")
	}
	if err := os.MkdirAll(cfg.LocalPath, os.ModePerm); err != nil {
		return nil, err
	}

	fileName := filepath.Base(cfg.URL)
	targetPath := filepath.Join(cfg.LocalPath, fileName)

	resp, err := http.Get(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return nil, err
	}

	return scanFiles(cfg.LocalPath)
}
