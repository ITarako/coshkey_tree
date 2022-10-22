package tree

import (
	"context"

	"github.com/ITarako/coshkey_tree/internal/model"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) GetFolder(ctx context.Context, id int64) (*model.Folder, error) {
	return s.repository.GetFolder(ctx, id)
}
