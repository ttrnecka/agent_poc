package repository

import (
	"context"

	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollectorRepository interface {
	GetByField(context.Context, string, interface{}) (*entity.Collector, error)
	GetByID(context.Context, primitive.ObjectID) (*entity.Collector, error)
	All(context.Context) ([]entity.Collector, error)
}

type collectorRepository struct {
	*cdb.CRUD[entity.Collector]
}

func NewCollectorRepository(db *cdb.CRUD[entity.Collector]) CollectorRepository {
	return &collectorRepository{db}
}
