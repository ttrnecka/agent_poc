package main

import (
	"encoding/gob"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/ttrnecka/agent_poc/webapi/db"

	logging "github.com/ttrnecka/agent_poc/logger"
)

var logger zerolog.Logger

func init() {
	logger = logging.SetupLogger("webapi")
}

func main() {
	// Needed for storing structs in sessions
	gob.Register(db.User{})

	// db
	err := db.Init()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	srv := &http.Server{
		Addr: ":8888",
		// Handler: Router(),
		Handler: EchoRouter(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error().Err(err).Msg("")
	}
}
