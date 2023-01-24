package models

//easyjson:json
type Users []User

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	About    string `json:"about"`
}

type UserUpdate struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
