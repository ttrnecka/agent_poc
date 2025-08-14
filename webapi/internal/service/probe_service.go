package service

import (
	"context"

	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProbeService interface {
	All(context.Context) ([]entity.Probe, error)
	GetProbe(context.Context, string) (*entity.Probe, error)
	DeleteProbe(context.Context, string) error
	UpdateProbe(context.Context, *entity.Probe) (primitive.ObjectID, error)
	Collector(context.Context, *entity.Probe) (*entity.Collector, error)
	// GetByName(context.Context, string) (*entity.User, error)
}

type probeService struct {
	probeRepo     repository.ProbeRepository
	collectorRepo repository.CollectorRepository
}

func NewProbeService(p repository.ProbeRepository, c repository.CollectorRepository) ProbeService {
	return &probeService{
		probeRepo:     p,
		collectorRepo: c,
	}
}

func (s *probeService) All(ctx context.Context) ([]entity.Probe, error) {
	return s.probeRepo.All(ctx)
}

func (s *probeService) GetProbe(ctx context.Context, id string) (*entity.Probe, error) {
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.probeRepo.GetByID(ctx, idp)
}

func (s *probeService) DeleteProbe(ctx context.Context, id string) error {
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.probeRepo.HardDeleteByID(ctx, idp)
}

func (s *probeService) UpdateProbe(ctx context.Context, probe *entity.Probe) (primitive.ObjectID, error) {
	if probe.ID.IsZero() {
		return s.probeRepo.Create(ctx, probe)
	}
	return probe.ID, s.probeRepo.UpdateByID(ctx, probe.ID, probe)
}

func (s *probeService) Collector(ctx context.Context, probe *entity.Probe) (*entity.Collector, error) {
	return s.collectorRepo.GetByID(ctx, probe.CollectorID)
}
