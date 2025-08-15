package handler

import (
	"encoding/json"
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

func (h *CollectorHandler) DeleteCollector(c echo.Context) error {
	probe_id := c.Param("id")
	_, err := h.service.Get(c.Request().Context(), probe_id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}
	err = h.service.Delete(c.Request().Context(), probe_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}

func (h *CollectorHandler) CreateUpdateCollector(c echo.Context) error {
	var collDTO dto.CollectorDTO
	if err := json.NewDecoder(c.Request().Body).Decode(&collDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	coll := mapper.ToCollectorEntity(collDTO)

	id, err := h.service.Update(c.Request().Context(), coll.ID, &coll)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})

	}
	collTmp, err := h.service.Get(c.Request().Context(), id.Hex())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	collDTO = mapper.ToCollectorDTO(*collTmp)
	return c.JSON(http.StatusOK, collDTO)
}
