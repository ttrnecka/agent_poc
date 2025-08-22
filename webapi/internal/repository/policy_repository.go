package repository

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
)

type PolicyRepository interface {
	GenericRepository[entity.Policy]
}

func NewPolicyRepository(db *cdb.CRUD[entity.Policy]) PolicyRepository {
	return NewGenericRepository(db)
}
