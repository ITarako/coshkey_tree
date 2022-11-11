package folder

import (
	"context"

	"coshkey_tree/internal/model"
	"coshkey_tree/internal/service/right"
	"github.com/rs/zerolog/log"
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
	folders, err := s.repository.AdminFolders(ctx)
	if err != nil {
		log.Error().Err(err).Msg("AdminFolders getting error")
		return
	}

	for _, folder := range folders {
		folderChan <- folder
	}
}

func (s Service) UserFolderCommand(ctx context.Context, user *model.User, folderChan chan<- model.Folder) {
	rootFolder, err := s.repository.RootFolder(ctx, user.Id)
	if err != nil {
		log.Error().Err(err).Msg("RootFolder getting error")
		return
	}

	folderChan <- *rootFolder

	folders, err := s.repository.Children(ctx, rootFolder.Lft, rootFolder.Rgt)
	if err != nil {
		log.Error().Err(err).Msg("Children getting error")
		return
	}

	for _, folder := range folders {
		folderChan <- folder
	}
}

func (s Service) UserAvailFolderForWriteCommand(ctx context.Context, user *model.User, rs right.Service, folderChan chan<- model.Folder) {
	rights := rs.SumRights(right.RightRead, right.RightWrite)

	folders, err := s.repository.UserFoldersByRight(ctx, user.Id, rights)
	if err != nil {
		log.Error().Err(err).Msg("UserFoldersByRight getting error")
		return
	}

	for _, folder := range folders {
		folderChan <- folder
	}
}

func (s Service) UserShareFolderCommand(ctx context.Context, user *model.User, folderChan chan<- model.Folder) {
	folders, err := s.repository.UserShareFolder(ctx, user.Id)
	if err != nil {
		log.Error().Err(err).Msg("UserShareFolder getting error")
		return
	}

	for _, folder := range folders {
		folderChan <- folder
	}
}

func (s Service) FavoriteCommand(ctx context.Context, user *model.User) (map[int]model.Folder, error) {
	favorites, err := s.repository.UserFavorite(ctx, user.Id)
	if err != nil {
		log.Error().Err(err).Msg("UserFavorite getting error")
		return nil, err
	}

	return favorites, nil
}
