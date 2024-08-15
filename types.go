package shoppinglistserver

import (
	"time"
)

type List struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Creator User    `json:"creator"`
	Entries []Entry `json:"entries,omitempty"`
}

type Entry struct {
	Id        int    `json:"id"`
	ListId    int    `json:"listId"`
	Text      string `json:"text"`
	Category  string `json:"category"`
	Completed bool   `json:"completed"`
}

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash,omitempty"`
}

type ListMember struct {
	ListId int
	UserId int
}

type Session struct {
	Token     string
	UserId    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
