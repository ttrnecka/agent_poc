package main

import (
	"encoding/gob"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/ttrnecka/agent_poc/webapi/db"
	"github.com/ttrnecka/agent_poc/webapi/server"
	"github.com/ttrnecka/agent_poc/webapi/shared/dto"

	logging "github.com/ttrnecka/agent_poc/logger"
)

var logger zerolog.Logger

func init() {
	logger = logging.SetupLogger("http")
}

func main() {
	// Needed for storing structs in sessions
	gob.Register(dto.UserDTO{})

	// db
	err := db.Init()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	srv := &http.Server{
		Addr: ":8888",
		// Handler: Router(),
		Handler: server.Router(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error().Err(err).Msg("")
	}
}
