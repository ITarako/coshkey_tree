package model

type User struct {
	Id       int32  `db:"id"`
	Username string `db:"username"`
	IsAdmin  bool
}
