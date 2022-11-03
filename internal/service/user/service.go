package user

import (
	"context"

	"github.com/ITarako/coshkey_tree/internal/model"
	internalErrors "github.com/ITarako/coshkey_tree/internal/pkg/errors"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) GetUser(ctx context.Context, id int) (*model.User, error) {
	if id == 0 {
		return nil, internalErrors.ErrNotFound
	}

	return s.repository.GetUser(ctx, id)
}
