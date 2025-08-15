package repository

import cdb "github.com/ttrnecka/agent_poc/common/db"

type GenericRepository[T any] interface {
	cdb.CRUDer[T]
}

type genericRepository[T any] struct {
	*cdb.CRUD[T]
}

func NewGenericRepository[T any](db *cdb.CRUD[T]) GenericRepository[T] {
	return &genericRepository[T]{db}
}
