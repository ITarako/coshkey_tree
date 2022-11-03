package model

const FolderClassName = "common\\models\\Folder"

const (
	FavoriteTypeFolder int = 1
	FavoriteTypeKey    int = 2
)

const FolderIdForTree int = -1

const (
	PrivateFolder  = "private-folder"
	PrivateProject = "private-project"
	SharedProject  = "shared-project"
	SharedFolder   = "shared-folder"
)

type Folder struct {
	Id        int    `db:"id"`
	IdUser    int    `db:"id_user"`
	IdParent  int    `db:"id_parent"`
	Lft       int    `db:"lft"`
	Rgt       int    `db:"rgt"`
	Depth     int    `db:"depth"`
	Title     string `db:"title"`
	IsActive  bool   `db:"is_active"`
	IsProject bool   `db:"is_project"`
	Children  map[int]Folder
}

type FavoriteFolder struct {
	Folder
	IsFavorite        bool `db:"is_favorite"`
	CountFavoriteKeys int  `db:"count_favorite_keys"`
	Children          map[int]FavoriteFolder
}

func (m Folder) GetIdUser() int {
	return m.IdUser
}

func (m Folder) GetIsProject() bool {
	return m.IsProject
}

type Folded interface {
	GetIdUser() int
	GetIsProject() bool
}

func SetClassification[T Folded](roots map[int]T, user *User) map[string][]T {
	var privateProject, sharedProject, privateFolder, sharedFolder []T

	for _, root := range roots {
		if root.GetIdUser() == user.Id && root.GetIsProject() {
			privateProject = append(privateProject, root)
		} else if root.GetIdUser() != user.Id && root.GetIsProject() {
			sharedProject = append(sharedProject, root)
		} else if root.GetIdUser() == user.Id && !root.GetIsProject() {
			privateFolder = append(privateFolder, root)
		} else if root.GetIdUser() != user.Id && !root.GetIsProject() {
			sharedFolder = append(sharedFolder, root)
		}
	}

	return map[string][]T{
		PrivateProject: privateProject,
		SharedProject:  sharedProject,
		PrivateFolder:  privateFolder,
		SharedFolder:   sharedFolder,
	}
}
