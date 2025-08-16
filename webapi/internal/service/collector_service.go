package service

import (
	"context"

	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollectorService interface {
	GenericService[entity.Collector]
}

type collectorService struct {
	GenericService[entity.Collector]
	probeRepo repository.ProbeRepository
}

func NewCollectorService(r repository.CollectorRepository, pr repository.ProbeRepository) CollectorService {
	s := &collectorService{
		GenericService: NewGenericService(r),
		probeRepo:      pr,
	}
	s.RegisterDependencies(s.DeleteProbes)
	return s
}

func (s *collectorService) DeleteProbes(ctx context.Context, id primitive.ObjectID) error {
	s.probeRepo.DeleteBy(ctx, "collector_id", id)
	return nil
}
