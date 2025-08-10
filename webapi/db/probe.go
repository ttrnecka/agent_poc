package db

import (
	"context"
	"fmt"

	cdb "github.com/ttrnecka/agent_poc/common/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Probe struct {
	cdb.BaseModel `bson:",inline"`
	Collector     string `bson:"collector" json:"collector"`
	Policy        string `bson:"policy" json:"policy"`
	Version       string `bson:"version" json:"version"`
	Address       string `bson:"address" json:"address"`
	Port          int    `bson:"port" json:"port"`
	User          string `bson:"user" json:"user"`
	// TODO hide this in json when probe saveing is properly resolved
	Password string `bson:"password" json:"password"`
}

func Probes() *cdb.CRUD[Probe] {
	return cdb.NewCRUD[Probe](dB.database, "probes")
}

func GetProbes(ctx context.Context) ([]Probe, error) {
	collection := dB.database.Collection("probes")

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "collectors"},
			{"localField", "collector_id"},
			{"foreignField", "_id"},
			{"as", "collector_info"},
		}}},
		{{"$unwind", "$collector_info"}},
		{{"$project", bson.D{
			{"policy", 1},
			{"version", 1},
			{"address", 1},
			{"port", 1},
			{"user", 1},
			{"password", 1},
			{"collector", "$collector_info.name"},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var probes []Probe
	if err := cursor.All(ctx, &probes); err != nil {
		return nil, err
	}
	return probes, nil
}

func SaveProbes(probes []interface{}) error {
	c := dB.database.Collection("probes")

	// Delete all documents in the collection
	if _, err := c.DeleteMany(context.Background(), bson.D{}); err != nil {
		return err
	}

	for _, probe := range probes {
		err := saveProbe(context.Background(), probe.(map[string]any))
		if err != nil {
			return err
		}
	}
	return nil
}

func saveProbe(ctx context.Context, probe map[string]interface{}) error {
	probesColl := dB.database.Collection("probes")

	collectorName, ok := probe["collector"].(string)
	if !ok {
		return fmt.Errorf("probe missing collector name")
	}
	coll, err := Collectors().GetByField(ctx, "name", collectorName)
	if err != nil {
		return fmt.Errorf("collector not found: %w", err)
	}

	// Replace collector name with collector_id
	delete(probe, "collector")
	probe["collector_id"] = coll.ID

	// Insert the probe
	_, err = probesColl.InsertOne(ctx, probe)
	if err != nil {
		return fmt.Errorf("failed to insert probe: %w", err)
	}
	return nil
}
