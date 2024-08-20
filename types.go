package shoppinglistserver

import (
	"time"
)

type List struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Creator User    `json:"creator,omitempty"`
	Entries []Entry `json:"entries,omitempty"`
	Members []User  `json:"members,omitempty"`
}

type Entry struct {
	Id        int       `json:"id"`
	ListId    int       `json:"listId"`
	Text      string    `json:"text"`
	Category  string    `json:"category"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
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

type Invitation struct {
	Token   string `json:"token"`
	Inviter User   `json:"inviter"`
	Invitee User   `json:"invitee"`
	List    List   `json:"list"`
}

type Session struct {
	Token     string    `json:"token"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

type Response struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
