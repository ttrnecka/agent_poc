package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type echoMiddleware struct {
}

func (e echoMiddleware) sessionManager() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			sess, err := session.Get(SESSION_STORE, c)
			if err != nil {
				return c.String(http.StatusInternalServerError, "failed to get session")
			}
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   86400 * 1,
				HttpOnly: true,
			}
			c.Set(SESSION_STORE, sess)
			return next(c)
		}
	}
}

func (e echoMiddleware) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := Session(c)
		logger.Info().Msg("AUTHENTICATING")
		if sess.Values["user"] != nil {
			return next(c)
		}
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}
