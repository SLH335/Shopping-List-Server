package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	. "github.com/slh335/shoppinglistserver"
)

func (server *Server) Register(c echo.Context) error {
	values, success, err := getFormValues(c, "username", "password")
	if !success {
		return err
	}
	username, password := values[0], values[1]

	user, err := server.AuthService.Register(username, password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to register user",
		})
	}

	session, err := server.AuthService.NewSession(user, 7)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to start new session",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "successfully registered user",
		Data:    session,
	})
}

func (server *Server) Login(c echo.Context) error {
	values, success, err := getFormValues(c, "username", "password")
	if !success {
		return err
	}
	username, password := values[0], values[1]

	user, err := server.AuthService.Login(username, password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: incorrect credentials",
		})
	}

	session, err := server.AuthService.NewSession(user, 7)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to start new session",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    session,
	})
}

func (server *Server) VerifySession(c echo.Context) error {
	values, success, err := getFormValues(c, "token")
	if !success {
		return err
	}
	token := values[0]

	session, err := server.AuthService.VerifySession(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: token invalid",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    session,
	})
}
