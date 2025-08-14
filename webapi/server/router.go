package server

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ttrnecka/agent_poc/webapi/api"
	"github.com/ttrnecka/agent_poc/webapi/internal/handler"
	mid "github.com/ttrnecka/agent_poc/webapi/server/middleware"
)

func Router() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(mid.SessionManager())

	userHandler := handler.NewUserHandler()

	e.POST("/api/login", userHandler.LoginUser)
	e.GET("/api/logout", userHandler.LogoutUser)
	e.GET("/api/user", userHandler.User, mid.AuthMiddleware)

	wsHandler := handler.NewWsHandler()
	e.GET("/ws", wsHandler.WS())

	ahandler := api.NewApiHandler()

	api := e.Group("/api/v1", mid.AuthMiddleware)
	api.GET("/collector", ahandler.CollectorsApiHandler)
	api.GET("/policy", ahandler.PoliciesApiHandler)
	api.GET("/probe", ahandler.ProbesApiHandler)
	api.POST("/probe", ahandler.ProbeCreateUpdateApiHandler)
	api.POST("/probe/:id", ahandler.ProbeCreateUpdateApiHandler)
	api.DELETE("/probe/:id", ahandler.ProbeDeleteApiHandler)

	api.GET("/policy/:name/:version", ahandler.PolicyApiHandler)

	api.GET("/data/collector", ahandler.DataApiCollectorsHandler)
	api.GET("/data/collector/:collector", ahandler.DataApiCollectorHandler)
	api.GET("/data/collector/:collector/:device", ahandler.DataApiCollectorDeviceHandler)
	api.GET("/data/collector/:collector/:device/:endpoint", ahandler.DataApiCollectorDeviceEndpointHandler)

	return e

}
