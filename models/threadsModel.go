package models

import "time"

//easyjson:json
type Threads []Thread

type Thread struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int32     `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}

type ThreadUpdate struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

const (
	Dislike = iota - 1
	Like    = iota
)

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int32  `json:"voice"`
}
