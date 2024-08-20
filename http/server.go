package http

import "github.com/slh335/shoppinglistserver/sqlite"

type Server struct {
	AuthService       *sqlite.AuthService
	UserService       *sqlite.UserService
	ListService       *sqlite.ListService
	EntryService      *sqlite.EntryService
	InvitationService *sqlite.InvitationService
}
