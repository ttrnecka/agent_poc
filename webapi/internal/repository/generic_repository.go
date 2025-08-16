package repository

import (
	"context"

	cdb "github.com/ttrnecka/agent_poc/common/db"
	"go.mongodb.org/mongo-driver/bson"
)

type GenericRepository[T any] interface {
	cdb.CRUDer[T]
	DeleteBy(context.Context, string, any) error
}

type genericRepository[T any] struct {
	*cdb.CRUD[T]
}

func NewGenericRepository[T any](db *cdb.CRUD[T]) GenericRepository[T] {
	return &genericRepository[T]{db}
}

func (r *genericRepository[T]) DeleteBy(ctx context.Context, key string, value any) error {
	return r.HardDelete(ctx, bson.M{key: value})
}
