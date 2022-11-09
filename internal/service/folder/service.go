package folder

import (
	"context"

	"coshkey_tree/internal/model"
	"coshkey_tree/internal/service/right"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) AdminFolderCommand(ctx context.Context, folderChan chan<- model.Folder) {
	folders, _ := s.repository.AdminFolders(ctx)
	for _, folder := range folders {
		folderChan <- folder
	}
}

func (s Service) UserFolderCommand(ctx context.Context, user *model.User, folderChan chan<- model.Folder) {
	rootFolder, err := s.repository.RootFolder(ctx, user.Id)
	if err != nil {
		return
	}
	folderChan <- *rootFolder

	folders, _ := s.repository.Children(ctx, rootFolder.Lft, rootFolder.Rgt)
	for _, folder := range folders {
		folderChan <- folder
	}
}

func (s Service) UserAvailFolderForWriteCommand(ctx context.Context, user *model.User, rs right.Service, folderChan chan<- model.Folder) {
	rights := rs.SumRights(right.RightRead, right.RightWrite)

	folders, _ := s.repository.UserFoldersByRight(ctx, user.Id, rights)
	for _, folder := range folders {
		folderChan <- folder
	}
}

func (s Service) UserShareFolderCommand(ctx context.Context, user *model.User, folderChan chan<- model.Folder) {
	folders, _ := s.repository.UserShareFolder(ctx, user.Id)
	for _, folder := range folders {
		folderChan <- folder
	}
}

func (s Service) FavoriteCommand(ctx context.Context, user *model.User) map[int]model.Folder {
	favorites, _ := s.repository.UserFavorite(ctx, user.Id)

	return favorites
}
