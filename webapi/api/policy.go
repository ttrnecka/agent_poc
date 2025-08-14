package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/agent_poc/webapi/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Policy struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

func (h *ApiHandler) PoliciesApiHandler(c echo.Context) error {
	policies, err := db.Policies().Find(c.Request().Context(), bson.D{})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, policies)
}

func (h *ApiHandler) PolicyApiHandler(c echo.Context) error {
	file := fmt.Sprintf("policies/bin/%s_%s", c.Param("name"), c.Param("version"))
	return c.File(file)
}
