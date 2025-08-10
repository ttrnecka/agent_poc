package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client   *mongo.Client
	database *mongo.Database
}

type Collector struct {
	Key  string                 `bson:"key" json:"key"`
	Data map[string]interface{} `bson:"data" json:"data"`
}

func Connect() (*DB, error) {
	uri := os.Getenv("MONGO_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db := DB{}
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db.database = client.Database("poc")
	db.client = client
	return &db, nil
}

func (db *DB) Database() *mongo.Database {
	return db.database
}

func (db *DB) Insert(collection string, data map[string]any) (*mongo.InsertOneResult, error) {
	c := db.client.Database("poc").Collection(collection)
	return c.InsertOne(context.Background(), data)
}

func (db *DB) Find(collection string, filter bson.D) (bson.M, error) {
	var found bson.M
	c := db.client.Database("poc").Collection(collection)
	err := c.FindOne(context.Background(), filter).Decode(&found)
	if err != nil {
		return nil, err
	}
	return found, nil
}

func (db *DB) FindAll(collection string, filter bson.D) ([]bson.M, error) {
	var results []bson.M
	c := db.client.Database("poc").Collection(collection)
	cur, err := c.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Fetch all collectors
func (db *DB) GetAllCollectors(ctx context.Context) ([]Collector, error) {
	colls, err := db.FindAll("collectors", bson.D{})
	if err != nil {
		return nil, err
	}
	var collectors []Collector
	for _, doc := range colls {
		// Marshal the bson.M into bytes
		data, err := bson.Marshal(doc)
		if err != nil {
			return nil, err
		}

		// Unmarshal bytes into Collector struct
		var c Collector
		if err := bson.Unmarshal(data, &c); err != nil {
			return nil, err
		}

		collectors = append(collectors, c)
	}
	return collectors, nil
}

// Fetch all collectors
func (db *DB) GetAllProbes(ctx context.Context) ([]map[string]any, error) {
	probes, err := db.FindAll("probes", bson.D{})
	if err != nil {
		return nil, err
	}
	var result []map[string]any = make([]map[string]any, len(probes))
	for i, doc := range probes {
		result[i] = map[string]any(doc)
	}
	return result, nil
}

// Fetch all collectors
func (db *DB) GetAllPolicies(ctx context.Context) ([]map[string]any, error) {
	p, err := db.FindAll("policies", bson.D{})
	if err != nil {
		return nil, err
	}
	var result []map[string]any = make([]map[string]any, len(p))
	for i, doc := range p {
		result[i] = map[string]any(doc)
	}
	return result, nil
}

func (db *DB) SaveProbes(probes []interface{}) error {
	c := db.client.Database("poc").Collection("probes")

	// Delete all documents in the collection
	if _, err := c.DeleteMany(context.Background(), bson.D{}); err != nil {
		return err
	}

	// Insert new documents
	if len(probes) > 0 {
		if _, err := c.InsertMany(context.Background(), probes); err != nil {
			return err
		}
	}
	return nil
}
