package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/agent_poc/webapi/db"
	"github.com/ttrnecka/agent_poc/webapi/ws"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) ProbesApiHandler(c echo.Context) error {
	probes, err := db.Probes().CRUD().All(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, probes)
}

func (h *Handler) ProbeCreateUpdateApiHandler(c echo.Context) error {
	var probe db.Probe
	if err := json.NewDecoder(c.Request().Body).Decode(&probe); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	id, err := probe.UpdateProbe(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})

	}
	probeTmp, err := db.Probes().CRUD().GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	go h.refreshPolicies()
	return c.JSON(http.StatusOK, probeTmp)
}

func (h *Handler) ProbeDeleteApiHandler(c echo.Context) error {
	probe_id := c.Param("id")
	id, err := primitive.ObjectIDFromHex(probe_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	probe, err := db.Probes().CRUD().GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	err = probe.Delete(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) refreshPolicies() {
	probes, err := db.Probes().CRUD().All(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Refreshing policies")
		return
	}

	collectors := make(map[string]bool)

	for _, p := range probes {
		collectors[p.Collector().Name] = true
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
