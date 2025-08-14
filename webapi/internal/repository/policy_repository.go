package repository

import (
	"context"

	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
)

type PolicyRepository interface {
	GetByField(context.Context, string, interface{}) (*entity.Policy, error)
	All(context.Context) ([]entity.Policy, error)
}

type policyRepository struct {
	*cdb.CRUD[entity.Policy]
}

func NewPolicyRepository(db *cdb.CRUD[entity.Policy]) PolicyRepository {
	return &policyRepository{db}
}
