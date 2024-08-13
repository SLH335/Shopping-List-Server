package shoppinglistserver

import (
	"time"
)

type Entry struct {
	Id        int    `json:"id"`
	Text      string `json:"text"`
	Category  string `json:"category"`
	Completed bool   `json:"completed"`
}

type User struct {
	Id           int
	Username     string
	PasswordHash string
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
