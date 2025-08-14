package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/agent_poc/webapi/internal/mapper"
	"github.com/ttrnecka/agent_poc/webapi/internal/service"
	"github.com/ttrnecka/agent_poc/webapi/shared/dto"
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
	var collectorsDTO []dto.CollectorDTO
	for _, col := range collectors {
		collectorsDTO = append(collectorsDTO, mapper.ToCollectorDTO(col))
	}
	return c.JSON(http.StatusOK, collectorsDTO)
}
