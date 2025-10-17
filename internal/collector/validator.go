package collector

import (
	"errors"
	"net/url"
	"strings"
)

// ValidateSource ensures the input source configuration is safe and valid.
func ValidateSource(cfg SourceConfig) error {
	if cfg.URL == "" {
		return errors.New("missing source URL")
	}

	if _, err := url.ParseRequestURI(cfg.URL); err != nil {
		return errors.New("invalid URL format")
	}

	switch cfg.Type {
	case "git":
		if !strings.Contains(cfg.URL, "github.com") &&
			!strings.Contains(cfg.URL, "gitlab.com") &&
			!strings.Contains(cfg.URL, "bitbucket.org") {
			return errors.New("unsupported git provider")
		}
	case "archive":
		if !strings.HasSuffix(cfg.URL, ".zip") &&
			!strings.HasSuffix(cfg.URL, ".tar.gz") &&
			!strings.HasSuffix(cfg.URL, ".tgz") {
			return errors.New("invalid archive format")
		}
	case "url":
		if !strings.HasPrefix(cfg.URL, "http://") &&
			!strings.HasPrefix(cfg.URL, "https://") {
			return errors.New("invalid http source")
		}
	default:
		return ErrInvalidSourceType
	}
	return nil
}
