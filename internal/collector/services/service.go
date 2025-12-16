package services

import (
	"github.com/unarya/unarya/internal/collector/domains"
	"github.com/unarya/unarya/pkg/types"
)

type CollectorService struct{}

func New() *CollectorService {
	return &CollectorService{}
}

func (s *CollectorService) CollectFromGit(cfg types.SourceConfig) (*types.CollectionResult, error) {
	return domains.CollectFromGit(cfg)
}

func (s *CollectorService) CollectFromArchive(cfg types.SourceConfig) (*types.CollectionResult, error) {
	return domains.CollectFromArchive(cfg)
}

func (s *CollectorService) CollectFromURL(cfg types.SourceConfig) (*types.CollectionResult, error) {
	return domains.CollectFromURL(cfg)
}

func (s *CollectorService) ValidateSource(cfg types.SourceConfig) error {
	return domains.ValidateSource(cfg)
}
