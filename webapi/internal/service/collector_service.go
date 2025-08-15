package service

import (
	"context"

	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollectorService interface {
	All(context.Context) ([]entity.Collector, error)
	Get(context.Context, string) (*entity.Collector, error)
	Delete(context.Context, string) error
	Update(context.Context, *entity.Collector) (primitive.ObjectID, error)
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

func (s *collectorService) Get(ctx context.Context, id string) (*entity.Collector, error) {
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, idp)
}

func (s *collectorService) Delete(ctx context.Context, id string) error {
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.HardDeleteByID(ctx, idp)
}

func (s *collectorService) Update(ctx context.Context, item *entity.Collector) (primitive.ObjectID, error) {
	if item.ID.IsZero() {
		return s.repo.Create(ctx, item)
	}
	return item.ID, s.repo.UpdateByID(ctx, item.ID, item)
}
