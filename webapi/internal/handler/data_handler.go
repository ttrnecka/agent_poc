package handler

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/agent_poc/webapi/internal/service"
)

type DataHandler struct {
	service service.DataService
}

func NewDataHandler(s service.DataService) *DataHandler {
	return &DataHandler{s}
}

func (h *DataHandler) Collectors(c echo.Context) error {
	collectors, err := h.service.Collectors()
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, collectors)
}

func (h *DataHandler) CollectorDevices(c echo.Context) error {
	devices, err := h.service.CollectorDevices(c.Param("collector"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"devices": devices})
}

func (h *DataHandler) CollectorDeviceEndpoints(c echo.Context) error {
	endpoints, err := h.service.CollectorDeviceEndpoints(c.Param("collector"), c.Param("device"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, map[string][]string{"endpoints": endpoints})
}

func (h *DataHandler) CollectorDeviceEndpointData(c echo.Context) error {
	content, err := h.service.CollectorDeviceEndpointData(c.Param("collector"), c.Param("device"), c.Param("endpoint"))
	if err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return c.String(http.StatusOK, string(content))
}
