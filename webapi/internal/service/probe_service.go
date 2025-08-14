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
	repo  repository.ProbeRepository
	crepo repository.CollectorRepository
}

func NewProbeService(r repository.ProbeRepository) ProbeService {
	return &probeService{repo: r}
}

func (s *probeService) All(ctx context.Context) ([]entity.Probe, error) {
	return s.repo.All(ctx)
}

func (s *probeService) GetProbe(ctx context.Context, id string) (*entity.Probe, error) {
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, idp)
}

func (s *probeService) DeleteProbe(ctx context.Context, id string) error {
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.HardDeleteByID(ctx, idp)
}

func (s *probeService) UpdateProbe(ctx context.Context, probe *entity.Probe) (primitive.ObjectID, error) {
	if probe.ID.IsZero() {
		return s.repo.Create(ctx, probe)
	}
	return probe.ID, s.repo.UpdateByID(ctx, probe.ID, probe)
}

func (s *probeService) Collector(ctx context.Context, probe *entity.Probe) (*entity.Collector, error) {
	if s.crepo == nil {
		s.crepo = repository.NewCollectorRepository(entity.Collectors())
	}
	return s.crepo.GetByID(ctx, probe.CollectorID)
}
