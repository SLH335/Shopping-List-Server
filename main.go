package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/slh335/einkaufsliste-server/models/sqlite"
)

type App struct {
	Entries *sqlite.EntryModel
}

func main() {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	app := App{
		Entries: &sqlite.EntryModel{
			DB: db,
		},
	}

	e := echo.New()

	e.GET("/entries", app.GetEntries)
	e.POST("/entry/complete", app.CompleteEntry)
	e.POST("/entry", app.InsertEntry)
	e.PUT("/entry", app.UpdateEntry)
	e.DELETE("/entry", app.DeleteEntry)

	e.Logger.Fatal(e.Start(":9001"))
}
