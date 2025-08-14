package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/agent_poc/webapi/internal/mapper"
	"github.com/ttrnecka/agent_poc/webapi/internal/service"
	"github.com/ttrnecka/agent_poc/webapi/server/ws"
	"github.com/ttrnecka/agent_poc/webapi/shared/dto"
)

type ProbeHandler struct {
	service service.ProbeService
}

func NewProbeHandler(s service.ProbeService) *ProbeHandler {
	return &ProbeHandler{s}
}

func (h *ProbeHandler) Probes(c echo.Context) error {

	probes, err := h.service.All(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, probes)
}

func (h *ProbeHandler) DeleteProbe(c echo.Context) error {
	probe_id := c.Param("id")
	_, err := h.service.GetProbe(c.Request().Context(), probe_id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}
	err = h.service.DeleteProbe(c.Request().Context(), probe_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}

func (h *ProbeHandler) CreateUpdateProbe(c echo.Context) error {
	var probeDTO dto.ProbeDTO
	if err := json.NewDecoder(c.Request().Body).Decode(&probeDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	probe := mapper.ToProbeEntity(probeDTO)

	id, err := h.service.UpdateProbe(c.Request().Context(), &probe)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})

	}
	probeTmp, err := h.service.GetProbe(c.Request().Context(), id.String())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	go h.refreshPolicies()
	probeDTO = mapper.ToProbeDTO(*probeTmp)
	return c.JSON(http.StatusOK, probeDTO)
}

func (h *ProbeHandler) refreshPolicies() {
	probes, err := h.service.All(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Refreshing policies")
		return
	}

	collectors := make(map[string]bool)

	for _, p := range probes {
		c, err := h.service.Collector(context.Background(), &p)
		if err != nil {
			logger.Error().Err(err).Msgf("Cannot find collector: %s", p.CollectorID)
		}
		collectors[c.Name] = true
	}

	// once saved we broadcast the new probes to all connected clients
	hub := ws.GetHub()
	for collector := range collectors {
		bmessage, err := json.Marshal(ws.NewMessage(ws.MSG_POLICY_REFRESH, "hub", collector, "Policy updated"))
		if err != nil {
			logger.Error().Err(err).Msg("Refreshing policies, marshall message")
			continue
		}
		hub.BroadcastMessage(bmessage)
	}
}
