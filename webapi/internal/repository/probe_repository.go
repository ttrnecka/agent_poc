package repository

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
)

type ProbeRepository interface {
	GenericRepository[entity.Probe]
}

func NewProbeRepository(db *cdb.CRUD[entity.Probe]) ProbeRepository {
	return NewGenericRepository(db)
}
