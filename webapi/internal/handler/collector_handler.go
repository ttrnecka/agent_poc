package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/agent_poc/webapi/internal/service"
)

type CollectorHandler struct {
	service service.CollectorService
}

func NewCollectorHandler(s service.CollectorService) *CollectorHandler {
	return &CollectorHandler{s}
}

func (h *CollectorHandler) Collectors(c echo.Context) error {

	collectors, err := h.service.All(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, collectors)
}
