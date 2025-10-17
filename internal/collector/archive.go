package collector

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// CollectFromArchive downloads and extracts a ZIP or TAR.GZ archive.
func CollectFromArchive(cfg SourceConfig) (*CollectionResult, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("archive url is empty")
	}
	if cfg.LocalPath == "" {
		cfg.LocalPath = filepath.Join(os.TempDir(), "collector_archive")
	}

	tmpFile := filepath.Join(os.TempDir(), filepath.Base(cfg.URL))
	defer os.Remove(tmpFile)

	resp, err := http.Get(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to download archive: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(tmpFile)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return nil, err
	}

	if strings.HasSuffix(cfg.URL, ".zip") {
		if err := extractZip(tmpFile, cfg.LocalPath); err != nil {
			return nil, err
		}
	} else if strings.HasSuffix(cfg.URL, ".tar.gz") || strings.HasSuffix(cfg.URL, ".tgz") {
		if err := extractTarGz(tmpFile, cfg.LocalPath); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unsupported archive format")
	}

	return scanFiles(cfg.LocalPath)
}

// extractZip extracts ZIP archives.
func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		out, err := os.Create(fpath)
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			out.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// extractTarGz extracts TAR.GZ archives.
func extractTarGz(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, os.ModePerm)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), os.ModePerm)
			out, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		}
	}
	return nil
}
