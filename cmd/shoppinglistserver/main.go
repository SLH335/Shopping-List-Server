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
		UserService: &sqlite.UserService{
			DB: db,
		},
		ListService: &sqlite.ListService{
			DB: db,
		},
		InvitationService: &sqlite.InvitationService{
			DB: db,
		},
		EntryService: &sqlite.EntryService{
			DB: db,
		},
	}

	e := echo.New()

	e.POST("/auth/register", server.Register)
	e.POST("/auth/login", server.Login)

	e.GET("/lists", server.GetLists)
	e.POST("/list", server.AddList)
	e.GET("/list/:id", server.GetEntries)
	e.DELETE("/list/:id", server.DeleteList)
	e.POST("/list/:id/join", server.JoinList)
	e.POST("/list/:id/leave", server.LeaveList)

	e.GET("/invitations", server.GetInvitations)
	e.POST("/invitation", server.Invite)
	e.DELETE("/invitation", server.RevokeInvitation)
	e.POST("/invitation/accept", server.AcceptInvitation)
	e.POST("/invitation/decline", server.DeclineInvitation)

	e.POST("/entry", server.AddEntry)
	e.PUT("/entry/:id", server.UpdateEntry)
	e.DELETE("/entry/:id", server.DeleteEntry)
	e.POST("/entry/:id/complete", server.CompleteEntry)

	e.Logger.Fatal(e.Start(":9001"))
}
