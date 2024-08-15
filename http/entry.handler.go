package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	. "github.com/slh335/shoppinglistserver"
)

func (server *Server) GetEntries(c echo.Context) error {
	_, success, err := verifySession(c, server)
	if !success {
		return err
	}

	listIdStr := c.Param("id")
	listId, err := strconv.Atoi(listIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: path parameter '%s' must be a valid integer", listIdStr),
		})
	}

	entries, err := server.EntryService.All(listId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load entries",
		})
	}
	return c.JSON(http.StatusOK, Response{Success: true, Data: entries})
}

func (server *Server) CompleteEntry(c echo.Context) error {
	_, success, err := verifySession(c, server)
	if !success {
		return err
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: path parameter '%s' must be a valid integer", idStr),
		})
	}

	values, success, err := getFormValues(c, "completed")
	if !success {
		return err
	}
	completedStr := values[0]
	completed := true
	if strings.ToLower(completedStr) == "false" {
		completed = false
	}

	updated, err := server.EntryService.Complete(id, completed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to complete entry",
		})
	}
	if !updated {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: entry %d does not exist", id),
		})
	}
	status := "complete"
	if !completed {
		status = "uncomplete"
	}
	entry, err := server.EntryService.Get(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load entry",
		})
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully marked entry %d as %s", id, status),
		Data:    entry,
	})
}

func (server *Server) AddEntry(c echo.Context) error {
	_, success, err := verifySession(c, server)
	if !success {
		return err
	}

	values, success, err := getFormValues(c, "list_id", "text", "category")
	if !success {
		return err
	}
	listIdStr, text, category := values[0], values[1], values[2]
	listId, err := strconv.Atoi(listIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: field 'list_id' must be a valid integer",
		})
	}

	entry, err := server.EntryService.Add(listId, text, category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to create entry",
		})
	}
	return c.JSON(http.StatusOK, Response{Success: true, Data: entry})
}

func (server *Server) UpdateEntry(c echo.Context) error {
	_, success, err := verifySession(c, server)
	if !success {
		return err
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: path parameter '%s' must be a valid integer", idStr),
		})
	}

	values, success, err := getFormValues(c, "text", "category")
	if !success {
		return err
	}
	text, category := values[0], values[1]

	updated, err := server.EntryService.Update(id, text, category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to update entry",
		})
	}
	if !updated {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: entry %d does not exist", id),
		})
	}
	entry, err := server.EntryService.Get(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load entry",
		})
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully updated entry %d", id),
		Data:    entry,
	})
}

func (server *Server) DeleteEntry(c echo.Context) error {
	_, success, err := verifySession(c, server)
	if !success {
		return err
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: path parameter '%s' must be a valid integer", idStr),
		})
	}

	deleted, err := server.EntryService.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to delete entry",
		})
	}
	if !deleted {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("error: entry %d does not exist", id),
		})
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully deleted entry %d", id),
	})
}
