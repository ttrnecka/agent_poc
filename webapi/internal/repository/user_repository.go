package repository

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
)

type UserRepository interface {
	cdb.CRUDer[entity.User]
}

func NewUserRepository(db *cdb.CRUD[entity.User]) UserRepository {
	return NewGenericRepository(db)
}
