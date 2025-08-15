package service

import (
	"context"

	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/internal/repository"
)

type ProbeService interface {
	GenericService[entity.Probe]
	Collector(context.Context, *entity.Probe) (*entity.Collector, error)
}

type probeService struct {
	GenericService[entity.Probe]
	collectorRepo repository.CollectorRepository
}

func NewProbeService(p repository.ProbeRepository, c repository.CollectorRepository) ProbeService {
	return &probeService{
		GenericService: NewGenericService(p),
		collectorRepo:  c,
	}
}

func (s *probeService) Collector(ctx context.Context, probe *entity.Probe) (*entity.Collector, error) {
	return s.collectorRepo.GetByID(ctx, probe.CollectorID)
}
