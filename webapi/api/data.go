package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

const baseDir = "/data/db"

func (h *ApiHandler) DataApiCollectorsHandler(c echo.Context) error {
	dirs, err := os.ReadDir(baseDir)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	var collectors []string
	for _, d := range dirs {
		if d.IsDir() {
			collectors = append(collectors, d.Name())
		}
	}
	return c.JSON(http.StatusOK, collectors)
}

// /api/v1/data/collector/:collector
func (h *ApiHandler) DataApiCollectorHandler(c echo.Context) error {
	collectorPath := filepath.Join(baseDir, c.Param("collector"))
	dirs, err := os.ReadDir(collectorPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	var devices []string
	for _, d := range dirs {
		if d.IsDir() {
			devices = append(devices, d.Name())
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"devices": devices})
}

func (h *ApiHandler) DataApiCollectorDeviceHandler(c echo.Context) error {
	devicePath := filepath.Join(baseDir, c.Param("collector"), c.Param("device"))

	entries, err := os.ReadDir(devicePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	var endpoints []string
	for _, entry := range entries {
		if !entry.IsDir() {
			endpoints = append(endpoints, entry.Name())
		}
	}
	return c.JSON(http.StatusOK, map[string][]string{"endpoints": endpoints})
}

func (h *ApiHandler) DataApiCollectorDeviceEndpointHandler(c echo.Context) error {
	filePath := filepath.Join(baseDir, c.Param("collector"), c.Param("device"), c.Param("endpoint"))

	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return c.String(http.StatusOK, string(content))
}
