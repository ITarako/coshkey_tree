package model

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	IsAdmin  bool
}
