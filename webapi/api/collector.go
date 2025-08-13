package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/agent_poc/webapi/db"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) CollectorsApiHandler(c echo.Context) error {

	collectors, err := db.Collectors().Find(c.Request().Context(), bson.D{})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, collectors)
}
