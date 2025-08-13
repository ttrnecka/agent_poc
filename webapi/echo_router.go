package main

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/api"
	"github.com/ttrnecka/agent_poc/webapi/db"
	"github.com/ttrnecka/agent_poc/webapi/ws"
	"golang.org/x/crypto/bcrypt"
)

type echoHandler struct {
}

const SESSION_STORE = "agentpoc"

func EchoRouter() *echo.Echo {
	e := echo.New()
	echoMiddleware := echoMiddleware{}

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(echoMiddleware.sessionManager())

	ehandler := echoHandler{}

	e.POST("/api/login", ehandler.login)
	e.GET("/api/logout", ehandler.logout)
	e.GET("/api/user", ehandler.user, echoMiddleware.authMiddleware)

	hub := ws.GetHub()
	go hub.Run()
	e.GET("/ws", ehandler.ws(hub))

	ahandler := api.NewHandler()

	api := e.Group("/api/v1", echoMiddleware.authMiddleware)
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

func (e echoHandler) ws(hub *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		return ws.ServeWs(hub, c.Response().Writer, c.Request())
	}
}

func (e echoHandler) login(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := db.Users().GetByField(c.Request().Context(), "username", username)

	if err != nil {
		logger.Error().Err(err).Msg("")
		if errors.Is(err, cdb.ErrNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Unexpected error: ", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	sess := Session(c)
	sess.Values["user"] = user
	if saveErr := sess.Save(c.Request(), c.Response()); saveErr != nil {
		return saveErr
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful", "user": user.Username})
}

func (e echoHandler) logout(c echo.Context) error {
	sess := Session(c)
	sess.Options.MaxAge = -1
	if saveErr := sess.Save(c.Request(), c.Response()); saveErr != nil {
		return saveErr
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out"})
}

func (e echoHandler) user(c echo.Context) error {
	sess := c.Get(SESSION_STORE).(*sessions.Session)
	user, ok := sess.Values["user"].(db.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Incorrect user in session, assertion failed")
	}
	return c.JSON(http.StatusOK, user)
}
