package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	ctx := context.Background()

	// Start MongoDB container
	req := tc.ContainerRequest{
		Image:        "mongo:7",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(20 * time.Second),
	}
	mongoC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := mongoC.Host(ctx)
	require.NoError(t, err)

	port, err := mongoC.MappedPort(ctx, "27017")
	require.NoError(t, err)

	uri := "mongodb://" + host + ":" + port.Port()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)

	cleanup := func() {
		_ = client.Disconnect(ctx)
		_ = mongoC.Terminate(ctx)
	}

	return &DB{client: client}, cleanup
}

func TestDB_InsertAndFind(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Insert a document
	data := map[string]any{
		"name":  "Alice",
		"email": "alice@example.com",
	}
	insertResult, err := db.Insert("users", data)
	require.NoError(t, err)
	require.NotNil(t, insertResult.InsertedID)

	// Find the document
	filter := bson.D{{Key: "name", Value: "Alice"}}
	found, err := db.Find("users", filter)
	require.NoError(t, err)
	require.Equal(t, "Alice", found["name"])
	require.Equal(t, "alice@example.com", found["email"])
}

func TestDB_FindAll(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Insert multiple documents
	users := []map[string]any{
		{"name": "Alice", "email": "alice@example.com"},
		{"name": "Bob", "email": "bob@example.com"},
		{"name": "Charlie", "email": "charlie@example.com"},
	}
	for _, user := range users {
		_, err := db.Insert("users", user)
		require.NoError(t, err)
	}

	// Find all documents with no filter
	all, err := db.FindAll("users", bson.D{})
	require.NoError(t, err)
	require.Len(t, all, 3)

	// Find only documents with name starting with "A" (regex)
	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: "^A"}}}}
	aUsers, err := db.FindAll("users", filter)
	require.NoError(t, err)
	require.Len(t, aUsers, 1)
	require.Equal(t, "Alice", aUsers[0]["name"])
}
