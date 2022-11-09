package tree

import (
	"context"
	"sync"

	"github.com/jmoiron/sqlx"

	"coshkey_tree/internal/model"
	"coshkey_tree/internal/service/folder"
	"coshkey_tree/internal/service/rbac"
	"coshkey_tree/internal/service/right"
	"coshkey_tree/internal/service/user"
)

type Result struct {
	Main     string `json:"main"`
	Favorite string `json:"favorite"`
}

type Service struct {
	db            *sqlx.DB
	UserService   user.Service
	FolderService folder.Service
	RbacService   rbac.Service
	RightService  right.Service
	CoshkeyUrl    string
}

func NewService(db *sqlx.DB, coshkeyUrl string) Service {
	return Service{
		UserService:   user.NewService(user.NewRepository(db)),
		FolderService: folder.NewService(folder.NewRepository(db)),
		RbacService:   rbac.NewService(rbac.NewRepository(db)),
		RightService:  right.NewService(),
		CoshkeyUrl:    coshkeyUrl,
	}
}

func (s Service) Generate(ctx context.Context, user *model.User, folderId int, ownTree bool) Result {
	res := Result{}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		res.Main = s.generateMain(ctx, user, folderId, ownTree)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		res.Favorite = s.generateFavorite(ctx, user, folderId)
	}()

	wg.Wait()

	return res
}

func (s Service) generateMain(ctx context.Context, user *model.User, folderId int, ownTree bool) string {
	var folderChan = make(chan model.Folder)

	if user.IsAdmin && !ownTree {
		go func() {
			defer close(folderChan)
			s.FolderService.AdminFolderCommand(ctx, folderChan)
		}()
	} else {
		go func() {
			defer close(folderChan)

			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				s.FolderService.UserFolderCommand(ctx, user, folderChan)
			}()

			if !ownTree {
				wg.Add(1)
				go func() {
					defer wg.Done()
					s.FolderService.UserShareFolderCommand(ctx, user, folderChan)
				}()
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				s.FolderService.UserAvailFolderForWriteCommand(ctx, user, s.RightService, folderChan)
			}()

			wg.Wait()
		}()

	}

	var folders = make(map[int]model.Folder)
	for f := range folderChan {
		folders[f.Id] = f
	}
	tree := s.buildMainTree(folders)
	return s.renderMain(tree, user, folderId)
}

func (s Service) generateFavorite(ctx context.Context, user *model.User, folderId int) string {
	folders := s.FolderService.FavoriteCommand(ctx, user)
	tree := s.buildFavoriteTree(folders)
	return s.renderFavorite(tree, user, folderId)
}

func (s Service) buildMainTree(folders map[int]model.Folder) map[int]model.Folder {
	rootsFolder := make(map[int]model.Folder)

	for id, f := range folders {
		parent, ok := folders[f.IdParent]
		if f.IdParent > 0 && ok {
			if len(parent.Children) == 0 {
				parent.Children = make(map[int]model.Folder)
			}

			parent.Children[id] = f
			folders[f.IdParent] = parent
			if parent.Depth == 1 {
				rootsFolder[parent.Id] = parent
			}
		}
		if f.Depth == 1 {
			rootsFolder[f.Id] = f
		}
	}

	return rootsFolder
}

func (s Service) buildFavoriteTree(folders map[int]model.Folder) map[int]model.Folder {
	rootsFolder := make(map[int]model.Folder)

	for id, f := range folders {
		parent, ok := folders[f.IdParent]
		if f.IdParent > 0 && ok {
			if len(parent.Children) == 0 {
				parent.Children = make(map[int]model.Folder)
			}

			parent.Children[id] = f
			folders[f.IdParent] = parent
			if parent.Depth == 1 {
				rootsFolder[parent.Id] = parent
			}
		}
		if f.Depth == 1 {
			rootsFolder[f.Id] = f
		}
	}

	roots := make(map[int]model.Folder)
	for id, f := range rootsFolder {
		if f.IsFavorite || f.CountFavoriteKeys > 0 || s.hasFavoriteChild(f) {
			roots[id] = f
		}
	}
	return roots
}

func (s Service) hasFavoriteChild(f model.Folder) bool {
	if len(f.Children) == 0 {
		return false
	}

	for _, child := range f.Children {
		if child.IsFavorite || child.CountFavoriteKeys > 0 || s.hasFavoriteChild(child) {
			return true
		}
	}

	return false
}
