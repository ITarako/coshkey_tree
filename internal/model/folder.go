package model

import (
	"sort"
	"strings"
)

const FolderClassName = "common\\models\\Folder"

const (
	FavoriteTypeFolder int = 1
	FavoriteTypeKey    int = 2
)

const FolderIdForTree int = -1

const (
	PrivateProject = 0
	SharedProject  = 1
	PrivateFolder  = 2
	SharedFolder   = 3
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

func (f Folder) GetLowerTitle() string {
	return strings.ToLower(f.Title)
}

func (f Folder) GetIdUser() int {
	return f.IdUser
}

func (f Folder) GetIsProject() bool {
	return f.IsProject
}

type Folded interface {
	GetLowerTitle() string
	GetIdUser() int
	GetIsProject() bool
}

func SetClassification[T Folded](roots map[int]T, user *User) [][]T {
	var privateProject, sharedProject, privateFolder, sharedFolder []T

	sortedSlice := GetSortedSliceFromMap(roots)

	for _, root := range sortedSlice {
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

	result := make([][]T, 4)
	result[PrivateProject] = privateProject
	result[SharedProject] = sharedProject
	result[PrivateFolder] = privateFolder
	result[SharedFolder] = sharedFolder
	return result
}

func GetSortedSliceFromMap[T Folded](folders map[int]T) []T {
	sortedSlice := make([]T, len(folders))
	i := 0
	for _, root := range folders {
		sortedSlice[i] = root
		i++
	}

	sort.Slice(sortedSlice, func(i, j int) bool {
		return sortedSlice[i].GetLowerTitle() < sortedSlice[j].GetLowerTitle()
	})

	return sortedSlice
}

func (f Folder) GetSortedChildren() []Folder {
	sortedSlice := make([]Folder, len(f.Children))
	i := 0
	for _, child := range f.Children {
		sortedSlice[i] = child
		i++
	}

	sort.Slice(sortedSlice, func(i, j int) bool {
		return sortedSlice[i].GetLowerTitle() < sortedSlice[j].GetLowerTitle()
	})

	return sortedSlice
}

func (f FavoriteFolder) GetSortedChildren() []FavoriteFolder {
	sortedSlice := make([]FavoriteFolder, len(f.Children))
	i := 0
	for _, child := range f.Children {
		sortedSlice[i] = child
		i++
	}

	sort.Slice(sortedSlice, func(i, j int) bool {
		return sortedSlice[i].GetLowerTitle() < sortedSlice[j].GetLowerTitle()
	})

	return sortedSlice
}
