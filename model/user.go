package model

type User struct {
	ID       int64  `json:"id"`
	Account  string `json:"account"`
	Password string `json:"password"`
}
