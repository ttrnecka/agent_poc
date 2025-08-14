package service

import (
	"context"

	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/internal/repository"
)

type PolicyService interface {
	All(context.Context) ([]entity.Policy, error)
	// GetByName(context.Context, string) (*entity.User, error)
}

type policyService struct {
	repo repository.PolicyRepository
}

func NewPolicyService(r repository.PolicyRepository) PolicyService {
	return &policyService{r}
}

func (s *policyService) All(ctx context.Context) ([]entity.Policy, error) {
	return s.repo.All(ctx)
}
