package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Probe struct {
	cdb.BaseModel `bson:",inline"`
	CollectorID   primitive.ObjectID `bson:"collector_id,omitempty"`
	Policy        string             `bson:"policy"`
	Version       string             `bson:"version"`
	Address       string             `bson:"address"`
	Port          int                `bson:"port"`
	User          string             `bson:"user"`
	Password      string             `bson:"password,omitempty"`
}

func Probes() *cdb.CRUD[Probe] {
	return cdb.NewCRUD[Probe](db.Database(), "probes")
}
