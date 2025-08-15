package service

import (
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/internal/repository"
)

type PolicyService interface {
	GenericService[entity.Policy]
}

type policyService struct {
	GenericService[entity.Policy]
}

func NewPolicyService(r repository.PolicyRepository) PolicyService {
	return &policyService{
		NewGenericService(r),
	}
}
