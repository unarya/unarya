package domain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/unarya/unarya/pkg/types"
)

// CollectFromGit clones a repository from a given Git URL and branch.
// Supports both public and private repositories (via token).
func CollectFromGit(cfg types.SourceConfig) (*types.CollectionResult, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("git url is empty")
	}
	if cfg.LocalPath == "" {
		cfg.LocalPath = filepath.Join(os.TempDir(), "collector_git_repo")
	}

	// If folder exists, remove and reclone
	os.RemoveAll(cfg.LocalPath)

	cloneURL := cfg.URL
	if cfg.Token != "" {
		// Inject token for authenticated clone
		cloneURL = injectToken(cfg.URL, cfg.Token)
	}

	cmdArgs := []string{"clone", "--depth", "1", "-b", cfg.Branch, cloneURL, cfg.LocalPath}
	if cfg.Branch == "" {
		cmdArgs = []string{"clone", "--depth", "1", cloneURL, cfg.LocalPath}
	}

	cmd := exec.Command("git", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git clone failed: %w", err)
	}

	// Walk collected files
	result, err := scanFiles(cfg.LocalPath)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// injectToken transforms https://github.com/user/repo.git to include a token.
func injectToken(url, token string) string {
	// Example: https://<token>@github.com/user/repo.git
	if token == "" {
		return url
	}
	return fmt.Sprintf("https://%s@%s", token, url[8:]) // strip https://
}
