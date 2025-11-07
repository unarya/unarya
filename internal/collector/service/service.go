package service

import (
	"github.com/unarya/unarya/internal/collector/domain"
	"github.com/unarya/unarya/pkg/types"
)

type CollectorService struct{}

func New() *CollectorService {
	return &CollectorService{}
}

func (s *CollectorService) CollectFromGit(cfg types.SourceConfig) (*types.CollectionResult, error) {
	return domain.CollectFromGit(cfg)
}

func (s *CollectorService) CollectFromArchive(cfg types.SourceConfig) (*types.CollectionResult, error) {
	return domain.CollectFromArchive(cfg)
}

func (s *CollectorService) CollectFromURL(cfg types.SourceConfig) (*types.CollectionResult, error) {
	return domain.CollectFromURL(cfg)
}

func (s *CollectorService) ValidateSource(cfg types.SourceConfig) error {
	return domain.ValidateSource(cfg)
}
