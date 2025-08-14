package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/agent_poc/webapi/db"
	"github.com/ttrnecka/agent_poc/webapi/server/middleware"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	// service service.UserService
}

// func NewUserHandler(s service.UserService) *UserHandler {
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (e *UserHandler) LoginUser(c echo.Context) error {

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

	sess := middleware.Session(c)
	sess.Values["user"] = user
	if saveErr := sess.Save(c.Request(), c.Response()); saveErr != nil {
		return saveErr
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful", "user": user.Username})
}

func (e *UserHandler) LogoutUser(c echo.Context) error {
	sess := middleware.Session(c)
	sess.Options.MaxAge = -1
	if saveErr := sess.Save(c.Request(), c.Response()); saveErr != nil {
		return saveErr
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out"})
}

func (e *UserHandler) User(c echo.Context) error {
	sess := middleware.Session(c)
	user, ok := sess.Values["user"].(db.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Incorrect user in session, assertion failed")
	}
	return c.JSON(http.StatusOK, user)
}

// func (h *UserHandler) CreateUser(c echo.Context) error {
// 	var req dto.CreateUserRequest
// 	if err := c.Bind(&req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	userEntity := mapper.ToUserEntity(req)
// 	createdUser, err := h.service.Create(userEntity)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusCreated, mapper.ToUserResponse(createdUser))
// }
