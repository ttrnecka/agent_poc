package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/agent_poc/webapi/internal/service"
)

type PolicyHandler struct {
	service service.PolicyService
}

func NewPolicyHandler(s service.PolicyService) *PolicyHandler {
	return &PolicyHandler{s}
}

func (h *PolicyHandler) Policies(c echo.Context) error {

	policies, err := h.service.All(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, policies)
}

func (h *PolicyHandler) PolicyFile(c echo.Context) error {
	file := fmt.Sprintf("policies/bin/%s_%s", c.Param("name"), c.Param("version"))
	return c.File(file)
}
