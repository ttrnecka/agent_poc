package db

import (
	"github.com/rs/zerolog"
	logging "github.com/ttrnecka/agent_poc/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

var logger zerolog.Logger
var dB *DB

func Database() *mongo.Database {
	return dB.database
}

func init() {
	logger = logging.SetupLogger("db")
}

func Init() error {
	d, err := Connect()
	if err != nil {
		return err
	}
	dB = d

	// Ensure all indexes before starting application logic
	if err := EnsureUserCollection(); err != nil {
		return err
	}
	return nil
}
