package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/agent_poc/webapi/server/ws"
)

type WsHandler struct {
	// service service.UserService
}

// func NewUserHandler(s service.UserService) *UserHandler {
func NewWsHandler() *WsHandler {
	return &WsHandler{}
}
func (h *WsHandler) WS() echo.HandlerFunc {
	hub := ws.GetHub()
	go hub.Run()
	return func(c echo.Context) error {
		return ws.ServeWs(hub, c.Response().Writer, c.Request())
	}
}
