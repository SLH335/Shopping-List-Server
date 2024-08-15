package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	. "github.com/slh335/shoppinglistserver"
)

func (server *Server) GetLists(c echo.Context) error {
	user, success, err := verifySession(c, server)
	if !success {
		return err
	}

	lists, err := server.ListService.All(user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load lists",
		})
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    lists,
	})
}

func (server *Server) AddList(c echo.Context) error {
	user, success, err := verifySession(c, server)
	if !success {
		return err
	}

	values, success, err := getFormValues(c, "name")
	if !success {
		return err
	}
	name := values[0]

	list, err := server.ListService.Add(user, name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to create list",
		})
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    list,
	})
}

func (server *Server) DeleteList(c echo.Context) error {
	user, success, err := verifySession(c, server)
	if !success {
		return err
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: path parameter '%s' is not a valid integer", idStr),
		})
	}

	list, err := server.ListService.Get(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load list",
		})
	}

	if list.Creator.Id != user.Id {
		return c.JSON(http.StatusForbidden, Response{
			Success: false,
			Message: "error: user is not authorized to delete that list",
		})
	}

	err = server.ListService.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to delete list",
		})
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "successfully deleted list",
	})
}

func (server *Server) JoinList(c echo.Context) error {
	user, success, err := verifySession(c, server)
	if !success {
		return err
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: path parameter '%s' is not a valid integer", idStr),
		})
	}

	err = server.ListService.Join(id, user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to join list",
		})
	}

	list, err := server.ListService.Get(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load list",
		})
	}
	entries, err := server.EntryService.All(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load entries",
		})
	}
	list.Entries = entries

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "successfully joined list",
		Data:    list,
	})
}

func (server *Server) LeaveList(c echo.Context) error {
	user, success, err := verifySession(c, server)
	if !success {
		return err
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: path parameter '%s' is not a valid integer", idStr),
		})
	}

	err = server.ListService.Leave(id, user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to leave list",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "successfully left list",
	})
}
