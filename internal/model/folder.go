package model

type Folder struct {
	Id        int32  `db:"id"`
	IdUser    int32  `db:"id_user"`
	IdParent  int32  `db:"id_parent"`
	Lft       int32  `db:"lft"`
	Rgt       int32  `db:"rgt"`
	Depth     int32  `db:"depth"`
	Title     string `db:"title"`
	IsActive  bool   `db:"is_active"`
	IsProject bool   `db:"is_project"`
}
