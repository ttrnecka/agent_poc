package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test model
type TestUser struct {
	BaseModel `bson:",inline"`
	Name      string `bson:"name"`
	Email     string `bson:"email"`
}

func (u *TestUser) SetCreatedUpdated() {
	u.BaseModel.SetCreatedUpdated()
}

func setupMongoContainer(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:7",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections"),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := mongoC.Host(ctx)
	require.NoError(t, err)

	port, err := mongoC.MappedPort(ctx, "27017")
	require.NoError(t, err)

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)

	db := client.Database("testdb")

	cleanup := func() {
		_ = client.Disconnect(ctx)
		_ = mongoC.Terminate(ctx)
	}

	return db, cleanup
}

func TestCRUDOperations(t *testing.T) {
	db, cleanup := setupMongoContainer(t)
	defer cleanup()

	userCRUD := NewCRUD[TestUser](db, "users")
	ctx := context.Background()

	// --- Create ---
	u := &TestUser{Name: "John", Email: "john@example.com"}
	id, err := userCRUD.Create(ctx, u)
	require.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, id)

	// --- GetByID ---
	found, err := userCRUD.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "John", found.Name)
	assert.Equal(t, "john@example.com", found.Email)

	// --- UpdateByID ---
	err = userCRUD.UpdateByID(ctx, id, bson.M{"email": "john.doe@example.com"})
	require.NoError(t, err)

	updated, err := userCRUD.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "john.doe@example.com", updated.Email)

	// --- Find ---
	results, err := userCRUD.Find(ctx, bson.M{"name": "John"})
	require.NoError(t, err)
	assert.Len(t, results, 1)

	// --- Pagination ---
	results, total, err := userCRUD.FindPaginated(ctx, bson.M{}, 1, 10, bson.D{{Key: "createdAt", Value: -1}})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, results, 1)

	// --- Soft Delete ---
	err = userCRUD.SoftDeleteByID(ctx, id)
	require.NoError(t, err)

	// Should not find soft-deleted doc
	_, err = userCRUD.GetByID(ctx, id)
	assert.ErrorIs(t, err, ErrNotFound)

	// --- Hard Delete ---
	err = userCRUD.HardDeleteByID(ctx, id)
	require.NoError(t, err)
}

func TestTimestamps(t *testing.T) {
	db, cleanup := setupMongoContainer(t)
	defer cleanup()

	userCRUD := NewCRUD[TestUser](db, "users")
	ctx := context.Background()

	u := &TestUser{Name: "Jane", Email: "jane@example.com"}
	id, err := userCRUD.Create(ctx, u)
	require.NoError(t, err)

	found, err := userCRUD.GetByID(ctx, id)
	require.NoError(t, err)

	assert.WithinDuration(t, time.Now(), found.CreatedAt, time.Second*2)
	assert.WithinDuration(t, time.Now(), found.UpdatedAt, time.Second*2)
}
