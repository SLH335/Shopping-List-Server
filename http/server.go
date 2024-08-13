package http

import "github.com/slh335/shoppinglistserver/sqlite"

type Server struct {
	AuthService  *sqlite.AuthService
	EntryService *sqlite.EntryService
}
