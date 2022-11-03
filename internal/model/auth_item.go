package model

const (
	TypeRole       int = 1
	TypePermission int = 2
)

type AuthItem struct {
	Name        string `db:"name"`
	Type        int    `db:"type"`
	Description string `db:"description"`
}
