package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/slh335/einkaufsliste-server/models"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (app *App) GetEntries(c echo.Context) error {
	entries, success, err := getAllEntries(app, c)
	if !success {
		return err
	}
	return c.JSON(http.StatusOK, Response{Success: true, Data: entries})
}

func (app *App) CompleteEntry(c echo.Context) error {
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
	entries, success, err := getAllEntries(app, c)
	if !success {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully marked entry %d as %s", id, status),
		Data:    entries,
	})
}

func (app *App) InsertEntry(c echo.Context) error {
	values, success, err := getFormValues(c, "text", "category")
	if !success {
		return err
	}
	text, category := values[0], values[1]

	_, err = app.Entries.Insert(text, category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to create entry",
		})
	}
	entries, success, err := getAllEntries(app, c)
	if !success {
		return err
	}
	return c.JSON(http.StatusOK, Response{Success: true, Data: entries})
}

func (app *App) UpdateEntry(c echo.Context) error {
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
	entries, success, err := getAllEntries(app, c)
	if !success {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully updated entry %d", id),
		Data:    entries,
	})
}

func (app *App) DeleteEntry(c echo.Context) error {
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
	entries, success, err := getAllEntries(app, c)
	if !success {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("successfully deleted entry %d", id),
		Data:    entries,
	})
}

func getFormValues(c echo.Context, keys ...string) (values []string, successs bool, err error) {
	values = []string{}
	for _, key := range keys {
		values = append(values, c.FormValue(key))
	}
	emptyValues := []string{}
	for _, value := range values {
		if value == "" {
			emptyValues = append(emptyValues, value)
		}
	}

	if len(emptyValues) > 0 {
		fields := ""
		if len(emptyValues) == 1 {
			fields += fmt.Sprintf("field '%s'", emptyValues[0])
		} else {
			fields += "fields"
			for i, emptyValue := range emptyValues {
				if i == len(emptyValue)-1 {
					fields += fmt.Sprintf("and '%s'", emptyValue)
				} else if 1 == len(emptyValue)-2 {
					fields += fmt.Sprintf(" '%s'", emptyValue)
				} else {
					fields += fmt.Sprintf(" '%s',", emptyValue)
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

func getAllEntries(app *App, c echo.Context) (entries []models.Entry, success bool, err error) {
	entries, err = app.Entries.All()
	if err != nil {
		err = c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load entries",
		})
		return entries, false, err
	}
	return entries, true, nil
}
