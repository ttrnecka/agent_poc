package server

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// db layer
	users := entity.Users()
	collectors := entity.Collectors()
	policies := entity.Policies()
	probes := entity.Probes()

	// repositories
	usersRepo := repository.NewUserRepository(users)
	collectorRepo := repository.NewCollectorRepository(collectors)
	policyRepo := repository.NewPolicyRepository(policies)
	probeRepo := repository.NewProbeRepository(probes)

	// services
	userSvc := service.NewUserService(usersRepo)
	collectorSvc := service.NewCollectorService(collectorRepo)
	policySvc := service.NewPolicyService(policyRepo)
	probeSvc := service.NewProbeService(probeRepo, collectorRepo)
	dataSvc := service.NewDataService("/data/db")

	//handlers
	userHandler := handler.NewUserHandler(userSvc)
	collectorHandler := handler.NewCollectorHandler(collectorSvc)
	policyHandler := handler.NewPolicyHandler(policySvc)
	probeHandler := handler.NewProbeHandler(probeSvc)
	wsHandler := handler.NewWsHandler()
	dataHandler := handler.NewDataHandler(dataSvc)

	e.POST("/api/login", userHandler.LoginUser)
	e.GET("/api/logout", userHandler.LogoutUser)
	e.GET("/api/user", userHandler.User, mid.AuthMiddleware)

	e.GET("/ws", wsHandler.WS())

	api := e.Group("/api/v1", mid.AuthMiddleware)
	api.GET("/collector", collectorHandler.Collectors)
	api.GET("/policy", policyHandler.Policies)
	api.GET("/probe", probeHandler.Probes)
	api.POST("/probe", probeHandler.CreateUpdateProbe)
	api.POST("/probe/:id", probeHandler.CreateUpdateProbe)
	api.DELETE("/probe/:id", probeHandler.DeleteProbe)

	api.GET("/policy/:name/:version", policyHandler.PolicyFile)

	api.GET("/data/collector", dataHandler.Collectors)
	api.GET("/data/collector/:collector", dataHandler.CollectorDevices)
	api.GET("/data/collector/:collector/:device", dataHandler.CollectorDeviceEndpoints)
	api.GET("/data/collector/:collector/:device/:endpoint", dataHandler.CollectorDeviceEndpointData)

	return e

}
