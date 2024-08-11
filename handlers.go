package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func (app *App) GetEntries(c echo.Context) error {
	entries, err := app.Entries.All()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(entries)
	return c.JSON(http.StatusOK, Response{Success: true, Data: entries})
}
