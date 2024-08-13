package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/slh335/shoppinglistserver/http"
	"github.com/slh335/shoppinglistserver/sqlite"
)

func main() {
	db, err := sqlite.Open("file:app.db")
	if err != nil {
		log.Fatal(err)
		return
	}

	server := http.Server{
		AuthService: &sqlite.AuthService{
			DB: db,
		},
		EntryService: &sqlite.EntryService{
			DB: db,
		},
	}

	e := echo.New()

	e.POST("/auth/register", server.Register)
	e.GET("/auth/login", server.Login)

	e.GET("/entries", server.GetEntries)
	e.POST("/entry/complete", server.CompleteEntry)
	e.POST("/entry", server.InsertEntry)
	e.PUT("/entry", server.UpdateEntry)
	e.DELETE("/entry", server.DeleteEntry)

	e.Logger.Fatal(e.Start(":9001"))
}
