package utils

import (
	"io"
	"os"
	"path/filepath"
)

// EnsureDir ensures directory exists
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// CopyFile copies file from src to dst
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
