package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("document not found")

// BaseModel for all documents
type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

type SetCreatedUpdateder interface {
	SetCreatedUpdated()
}

type CRUDer[T any] interface {
	GetByID(context.Context, primitive.ObjectID) (*T, error)
	GetByField(context.Context, string, any) (*T, error)
	All(context.Context) ([]T, error)
	HardDeleteByID(context.Context, primitive.ObjectID) error
	HardDelete(context.Context, interface{}, ...*options.DeleteOptions) error
	Create(context.Context, *T) (primitive.ObjectID, error)
	UpdateByID(context.Context, primitive.ObjectID, *T) error
	GetCollection() *mongo.Collection
	Find(context.Context, interface{}, ...*options.FindOptions) ([]T, error)
	FindWithSoftDeleted(context.Context, interface{}, ...*options.FindOptions) ([]T, error)
	InsertAll(context.Context, []T) error
	SoftDeleteByID(context.Context, primitive.ObjectID) error
	RestoreByID(context.Context, primitive.ObjectID) error
}

func (m *BaseModel) SetCreatedUpdated() {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	m.UpdatedAt = now
}

func (m *BaseModel) GetID() primitive.ObjectID {
	return m.ID
}

type CRUD[T any] struct {
	Collection *mongo.Collection
}

func NewCRUD[T any](db *mongo.Database, collectionName string) *CRUD[T] {
	return &CRUD[T]{Collection: db.Collection(collectionName)}
}

func (c *CRUD[T]) GetCollection() *mongo.Collection {
	return c.Collection
}

func (c *CRUD[T]) Create(ctx context.Context, doc *T) (primitive.ObjectID, error) {
	if bm, ok := any(doc).(interface{ SetCreatedUpdated() }); ok {
		bm.SetCreatedUpdated()
	}
	res, err := c.Collection.InsertOne(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (c *CRUD[T]) GetByID(ctx context.Context, id primitive.ObjectID) (*T, error) {
	var result T
	err := c.Collection.FindOne(ctx, bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	return &result, err
}

func (c *CRUD[T]) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	f := mergeFilters(filter, bson.M{"deletedAt": bson.M{"$exists": false}})
	cursor, err := c.Collection.Find(ctx, f, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	// return empty array if nothing is found
	if results == nil {
		return []T{}, nil
	}
	return results, nil
}

func (c *CRUD[T]) All(ctx context.Context) ([]T, error) {
	return c.Find(ctx, bson.D{})
}

func (c *CRUD[T]) FindWithSoftDeleted(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	cursor, err := c.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	// return empty array if nothing is found
	if results == nil {
		return []T{}, nil
	}
	return results, nil
}

func (c *CRUD[T]) GetByField(ctx context.Context, field string, value interface{}) (*T, error) {
	var result T
	err := c.Collection.FindOne(ctx, bson.M{
		field:       value,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	return &result, err
}

func (c *CRUD[T]) FindPaginated(ctx context.Context, filter interface{}, page, pageSize int64, sort bson.D) ([]T, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	f := mergeFilters(filter, bson.M{"deletedAt": bson.M{"$exists": false}})
	opts := options.Find().
		SetSkip((page - 1) * pageSize).
		SetLimit(pageSize).
		SetSort(sort)

	cursor, err := c.Collection.Find(ctx, f, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	count, err := c.Collection.CountDocuments(ctx, f)
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func (c *CRUD[T]) UpdateByID(ctx context.Context, id primitive.ObjectID, doc *T) error {
	update := bson.M{
		"$set":         doc,
		"$currentDate": bson.M{"updatedAt": true},
	}
	if bm, ok := any(doc).(SetCreatedUpdateder); ok {
		bm.SetCreatedUpdated()
		update = bson.M{
			"$set": doc,
		}
	}
	_, err := c.Collection.UpdateOne(ctx,
		bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}, update)
	return err
}

func (c *CRUD[T]) InsertAll(ctx context.Context, docs []T) error {
	inserts := make([]any, 0)
	for _, d := range docs {
		if bm, ok := any(d).(interface{ SetCreatedUpdated() }); ok {
			bm.SetCreatedUpdated()
		}
		inserts = append(inserts, d)
	}
	_, err := c.Collection.InsertMany(ctx, inserts)
	return err
}

func (c *CRUD[T]) SoftDeleteByID(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	_, err := c.Collection.UpdateOne(ctx,
		bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"deletedAt": now}})
	return err
}

func (c *CRUD[T]) RestoreByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := c.Collection.UpdateOne(ctx,
		bson.M{"_id": id, "deletedAt": bson.M{"$exists": true}},
		bson.M{"$unset": bson.M{"deletedAt": ""}})
	return err
}

func (c *CRUD[T]) HardDeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := c.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (c *CRUD[T]) HardDelete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) error {
	_, err := c.Collection.DeleteMany(ctx, filter, opts...)
	return err
}

func mergeFilters(userFilter interface{}, extraFilter bson.M) bson.M {
	m := bson.M{}
	if userFilter != nil {
		if um, ok := userFilter.(bson.M); ok {
			for k, v := range um {
				m[k] = v
			}
		}
	}
	for k, v := range extraFilter {
		m[k] = v
	}
	return m
}
