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
	Id                int    `db:"id"`
	IdUser            int    `db:"id_user"`
	IdParent          int    `db:"id_parent"`
	Lft               int    `db:"lft"`
	Rgt               int    `db:"rgt"`
	Depth             int    `db:"depth"`
	Title             string `db:"title"`
	IsActive          bool   `db:"is_active"`
	IsProject         bool   `db:"is_project"`
	IsFavorite        bool   `db:"is_favorite"`
	CountFavoriteKeys int    `db:"count_favorite_keys"`
	Children          map[int]Folder
}

func (f Folder) GetLowerTitle() string {
	return strings.ToLower(f.Title)
}

func SetClassification(roots map[int]Folder, user *User) [][]Folder {
	var privateProject, sharedProject, privateFolder, sharedFolder []Folder

	sortedSlice := GetSortedSliceFromMap(roots)

	for _, root := range sortedSlice {
		if root.IdUser == user.Id && root.IsProject {
			privateProject = append(privateProject, root)
		} else if root.IdUser != user.Id && root.IsProject {
			sharedProject = append(sharedProject, root)
		} else if root.IdUser == user.Id && !root.IsProject {
			privateFolder = append(privateFolder, root)
		} else if root.IdUser != user.Id && !root.IsProject {
			sharedFolder = append(sharedFolder, root)
		}
	}

	result := make([][]Folder, 4)
	result[PrivateProject] = privateProject
	result[SharedProject] = sharedProject
	result[PrivateFolder] = privateFolder
	result[SharedFolder] = sharedFolder
	return result
}

func GetSortedSliceFromMap(folders map[int]Folder) []Folder {
	sortedSlice := make([]Folder, len(folders))
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
