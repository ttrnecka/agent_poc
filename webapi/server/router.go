package server

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ttrnecka/agent_poc/webapi/api"
	"github.com/ttrnecka/agent_poc/webapi/internal/entity"
	"github.com/ttrnecka/agent_poc/webapi/internal/handler"
	"github.com/ttrnecka/agent_poc/webapi/internal/repository"
	"github.com/ttrnecka/agent_poc/webapi/internal/service"
	mid "github.com/ttrnecka/agent_poc/webapi/server/middleware"
)

func Router() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(mid.SessionManager())

	userHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(entity.Users())))

	e.POST("/api/login", userHandler.LoginUser)
	e.GET("/api/logout", userHandler.LogoutUser)
	e.GET("/api/user", userHandler.User, mid.AuthMiddleware)

	wsHandler := handler.NewWsHandler()
	e.GET("/ws", wsHandler.WS())

	ahandler := api.NewApiHandler()

	collectorHandler := handler.NewCollectorHandler(service.NewCollectorService(repository.NewCollectorRepository(entity.Collectors())))
	policyHandler := handler.NewPolicyHandler(service.NewPolicyService(repository.NewPolicyRepository(entity.Policies())))
	probeHandler := handler.NewProbeHandler(service.NewProbeService(repository.NewProbeRepository(entity.Probes())))

	api := e.Group("/api/v1", mid.AuthMiddleware)
	api.GET("/collector", collectorHandler.Collectors)
	api.GET("/policy", policyHandler.Policies)
	api.GET("/probe", probeHandler.Probes)
	api.POST("/probe", probeHandler.CreateUpdateProbe)
	api.POST("/probe/:id", probeHandler.CreateUpdateProbe)
	api.DELETE("/probe/:id", probeHandler.DeleteProbe)

	api.GET("/policy/:name/:version", policyHandler.PolicyFile)

	api.GET("/data/collector", ahandler.DataApiCollectorsHandler)
	api.GET("/data/collector/:collector", ahandler.DataApiCollectorHandler)
	api.GET("/data/collector/:collector/:device", ahandler.DataApiCollectorDeviceHandler)
	api.GET("/data/collector/:collector/:device/:endpoint", ahandler.DataApiCollectorDeviceEndpointHandler)

	return e

}
