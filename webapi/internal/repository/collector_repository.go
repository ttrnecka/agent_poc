package repository

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
)

type CollectorRepository interface {
	cdb.CRUDer[entity.Collector]
}

func NewCollectorRepository(db *cdb.CRUD[entity.Collector]) CollectorRepository {
	return NewGenericRepository(db)
}
