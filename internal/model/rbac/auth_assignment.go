package model

type AuthAssignment struct {
	ItemName string `db:"item_name"`
	UserId   int32  `db:"user_id"`
}
