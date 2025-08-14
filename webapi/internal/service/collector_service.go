package service

import (
	"context"

	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/internal/repository"
)

type CollectorService interface {
	All(context.Context) ([]entity.Collector, error)
	// GetByName(context.Context, string) (*entity.User, error)
}

type collectorService struct {
	repo repository.CollectorRepository
}

func NewCollectorService(r repository.CollectorRepository) CollectorService {
	return &collectorService{r}
}

func (s *collectorService) All(ctx context.Context) ([]entity.Collector, error) {
	return s.repo.All(ctx)
}
