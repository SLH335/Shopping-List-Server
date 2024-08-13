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
	success, err := verifySession(c, server)
	if !success {
		return err
	}

	entries, success, err := getAllEntries(server, c)
	if !success {
		return err
	}
	return c.JSON(http.StatusOK, Response{Success: true, Data: entries})
}

func (server *Server) CompleteEntry(c echo.Context) error {
	success, err := verifySession(c, server)
	if !success {
		return err
	}

	values, success, err := getFormValues(c, "id", "completed")
	if !success {
		return err
	}
	idStr, completedStr := values[0], values[1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: field 'id' must be a valid integer",
		})
	}
	completed := true
	if strings.ToLower(completedStr) == "false" {
		completed = false
	}

	updated, err := server.EntryService.Complete(id, completed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to create entry",
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
	entries, success, err := getAllEntries(server, c)
	if !success {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully marked entry %d as %s", id, status),
		Data:    entries,
	})
}

func (server *Server) InsertEntry(c echo.Context) error {
	success, err := verifySession(c, server)
	if !success {
		return err
	}

	values, success, err := getFormValues(c, "text", "category")
	if !success {
		return err
	}
	text, category := values[0], values[1]

	_, err = server.EntryService.Insert(text, category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to create entry",
		})
	}
	entries, success, err := getAllEntries(server, c)
	if !success {
		return err
	}
	return c.JSON(http.StatusOK, Response{Success: true, Data: entries})
}

func (server *Server) UpdateEntry(c echo.Context) error {
	success, err := verifySession(c, server)
	if !success {
		return err
	}

	values, success, err := getFormValues(c, "id", "text", "category")
	if !success {
		return err
	}
	idStr, text, category := values[0], values[1], values[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: field 'id' must be a valid integer",
		})
	}

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
	entries, success, err := getAllEntries(server, c)
	if !success {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully updated entry %d", id),
		Data:    entries,
	})
}

func (server *Server) DeleteEntry(c echo.Context) error {
	success, err := verifySession(c, server)
	if !success {
		return err
	}

	values, success, err := getFormValues(c, "id")
	if !success {
		return err
	}
	idStr := values[0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: field 'id' must be a valid integer",
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
	entries, success, err := getAllEntries(server, c)
	if !success {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully deleted entry %d", id),
		Data:    entries,
	})
}

func getAllEntries(server *Server, c echo.Context) (entries []Entry, success bool, err error) {
	entries, err = server.EntryService.All()
	if err != nil {
		err = c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load entries",
		})
		return entries, false, err
	}
	return entries, true, nil
}
