package model

const (
	TYPE_ROLE       int32 = 1
	TYPE_PERMISSION int32 = 2
)

type AuthItem struct {
	Name        string `db:"name"`
	Type        int32  `db:"type"`
	Description string `db:"description"`
}
