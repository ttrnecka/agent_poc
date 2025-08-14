package db

import (
	"context"

	cdb "github.com/ttrnecka/agent_poc/common/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type probes struct {
	crud *cdb.CRUD[Probe]
}

func Probes() probes {
	return probes{
		crud: cdb.NewCRUD[Probe](dB.database, "probes"),
	}
}

func (p probes) CRUD() *cdb.CRUD[Probe] {
	return p.crud
}

func (p probes) Probes(ctx context.Context) ([]Probe, error) {
	return p.crud.All(ctx)
}

type Probe struct {
	cdb.BaseModel `bson:",inline"`
	// Collector     string             `bson:"-" json:"collector"`
	CollectorID primitive.ObjectID `bson:"collector_id,omitempty" json:"collector_id"`
	Policy      string             `bson:"policy" json:"policy"`
	Version     string             `bson:"version" json:"version"`
	Address     string             `bson:"address" json:"address"`
	Port        int                `bson:"port" json:"port"`
	User        string             `bson:"user" json:"user"`
	// TODO hide this in json when probe saveing is properly resolved
	Password string `bson:"password" json:"password"`
}

// func (p *Probe) Collector() *db.Collector {
// 	// coll, _ := Collectors().GetByID(context.Background(), p.CollectorID)
// 	// return coll
// 	return nil
// }

func (p *Probe) UpdateProbe(ctx context.Context) (primitive.ObjectID, error) {
	crud := Probes().CRUD()
	if p.ID.IsZero() {
		return crud.Create(ctx, p)
	}
	return p.ID, crud.UpdateByID(ctx, p.ID, p)
}

func (p *Probe) Delete(ctx context.Context) error {
	return Probes().CRUD().HardDeleteByID(ctx, p.ID)
}
