package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	. "github.com/slh335/shoppinglistserver"
)

func getFormValues(c echo.Context, keys ...string) (values []string, successs bool, err error) {
	for _, key := range keys {
		values = append(values, c.FormValue(key))
	}

	emptyKeys := []string{}
	for _, key := range keys {
		if c.FormValue(key) == "" {
			emptyKeys = append(emptyKeys, key)
		}
	}

	if len(emptyKeys) > 0 {
		fields := ""
		if len(emptyKeys) == 1 {
			fields += fmt.Sprintf("field '%s'", emptyKeys[0])
		} else {
			fields += "fields"
			for i, emptyKey := range emptyKeys {
				if i == len(emptyKey)-1 {
					fields += fmt.Sprintf("and '%s'", emptyKey)
				} else if 1 == len(emptyKey)-2 {
					fields += fmt.Sprintf(" '%s'", emptyKey)
				} else {
					fields += fmt.Sprintf(" '%s',", emptyKey)
				}
			}
		}
		err = c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: %s must be provided", fields),
		})
		return values, false, err
	}
	return values, true, nil
}

func verifySession(c echo.Context, server *Server) (user User, success bool, err error) {
	sessionToken := strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", 1)

	session, err := server.AuthService.VerifySession(sessionToken)
	if err != nil {
		err = c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: invalid credentials",
		})
		return User{}, false, err
	}
	return session.User, true, nil
}
