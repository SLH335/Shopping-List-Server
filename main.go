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

	e.Logger.Fatal(e.Start(":9001"))
}
