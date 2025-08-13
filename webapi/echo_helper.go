package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func Session(c echo.Context) *sessions.Session {
	if sess, ok := c.Get(SESSION_STORE).(*sessions.Session); ok {
		return sess
	}
	return nil
}
