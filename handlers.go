package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (app *App) GetEntries(c echo.Context) error {
	entries, err := app.Entries.All()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load entries",
		})
	}
	return c.JSON(http.StatusOK, Response{Success: true, Data: entries})
}

func (app *App) CompleteEntry(c echo.Context) error {
	idStr := c.FormValue("id")
	completedStr := c.FormValue("completed")

	if idStr == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: field 'id' must be provided",
		})
	}
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

	updated, err := app.Entries.Complete(id, completed)
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
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully marked entry %d as %s", id, status),
		Data:    id,
	})
}

func (app *App) InsertEntry(c echo.Context) error {
	text := c.FormValue("text")
	category := c.FormValue("category")

	if text == "" || category == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: fields 'text' and 'category' must be provided",
		})
	}

	id, err := app.Entries.Insert(text, category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to create entry",
		})
	}
	return c.JSON(http.StatusOK, Response{Success: true, Data: id})
}

func (app *App) UpdateEntry(c echo.Context) error {
	idStr := c.FormValue("id")
	text := c.FormValue("text")
	category := c.FormValue("category")

	if idStr == "" || text == "" || category == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: fields 'id', 'text' and 'category' must be provided",
		})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: field 'id' must be a valid integer",
		})
	}

	updated, err := app.Entries.Update(id, text, category)
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
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully updated entry %d", id),
	})
}

func (app *App) DeleteEntry(c echo.Context) error {
	idStr := c.FormValue("id")

	if idStr == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: field 'id' must be provided",
		})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "error: field 'id' must be a valid integer",
		})
	}

	deleted, err := app.Entries.Delete(id)
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
