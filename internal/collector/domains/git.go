package domains

import (
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/unarya/unarya/internal/infrastructures"
	"github.com/unarya/unarya/pkg/types"
)

// CollectFromGit clones a repository from a given Git URL and branch.
// Supports both public and private repositories (via token).
func CollectFromGit(cfg types.SourceConfig) (*types.CollectionResult, error) {
	minioClient := infrastructures.MinioClient
	bucket := infrastructures.BucketName
	endpoint := infrastructures.Endpoint
	if cfg.URL == "" {
		return nil, fmt.Errorf("git url is empty")
	}
	if cfg.Path == "" {
		cfg.Path = filepath.Join(os.TempDir(), "collector_git_repo")
	}
	var rP string

	// Clean local directory
	if err := os.RemoveAll(cfg.Path); err != nil {
		return nil, err
	}

	// Prepare clone args
	args := []string{"clone", "--depth=1"}
	if cfg.Branch != "" {
		args = append(args, "-b", cfg.Branch)
	}
	args = append(args, cfg.URL, cfg.Path)

	// Git clone with context timeout
	cmd := exec.CommandContext(cfg.Ctx, "git", args...)

	// Set GIT_TERMINAL_PROMPT=0 to prevent git from asking for a password when cloning a private repo without a token
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "Authentication failed") ||
			strings.Contains(string(out), "could not read Username") {
			return nil, fmt.Errorf("authentication failed: this repository may be private and requires a valid token")
		}
		return nil, fmt.Errorf("git clone failed: %v\n%s", err, out)
	}

	// Upload all files to MinIO
	err = filepath.Walk(cfg.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// Skip .git directory
		if strings.Contains(path, ".git"+string(filepath.Separator)) {
			return nil
		}

		relPath, _ := filepath.Rel("github.com", path)
		rP = relPath
		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = minioClient.PutObject(
			cfg.Ctx,
			bucket,
			relPath,
			file,
			info.Size(),
			minio.PutObjectOptions{
				ContentType: contentType,
			},
		)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to upload to minio: %w", err)
	}

	// Build scan result
	result, err := scanFiles(cfg.Path)
	if err != nil {
		return nil, err
	}
	result.Dir = fmt.Sprintf("[S3 (%v): %s]", endpoint+bucket, rP)
	// Cleanup local directory after upload
	defer os.RemoveAll(cfg.Path)

	return result, nil
}
