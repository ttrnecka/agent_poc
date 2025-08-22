package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	logging "github.com/ttrnecka/agent_poc/logger"
)

var logger zerolog.Logger
var validate *validator.Validate

func init() {
	logger = logging.SetupLogger("handler")
	validate = validator.New(validator.WithRequiredStructEnabled())
}
