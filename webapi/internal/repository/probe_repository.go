package repository

import (
	"context"

	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProbeRepository interface {
	GetByField(context.Context, string, interface{}) (*entity.Probe, error)
	GetByID(context.Context, primitive.ObjectID) (*entity.Probe, error)
	All(context.Context) ([]entity.Probe, error)
	HardDeleteByID(context.Context, primitive.ObjectID) error
	Create(context.Context, *entity.Probe) (primitive.ObjectID, error)
	UpdateByID(context.Context, primitive.ObjectID, *entity.Probe) error
}

type probeRepository struct {
	*cdb.CRUD[entity.Probe]
}

func NewProbeRepository(db *cdb.CRUD[entity.Probe]) ProbeRepository {
	return &probeRepository{db}
}
