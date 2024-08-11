package models

type Entry struct {
	Id        int    `json:"id"`
	Text      string `json:"text"`
	Category  string `json:"category"`
	Completed bool   `json:"completed"`
}
