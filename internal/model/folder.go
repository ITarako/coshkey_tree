package model

type Folder struct {
	Id        int64  `db:"id"`
	IdUser    int64  `db:"id_user"`
	IdParent  int64  `db:"id_parent"`
	Lft       int64  `db:"lft"`
	Rgt       int64  `db:"rgt"`
	Depth     int64  `db:"depth"`
	Title     string `db:"title"`
	IsActive  bool   `db:"is_active"`
	IsProject bool   `db:"is_project"`
}
