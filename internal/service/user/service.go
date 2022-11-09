package user

import (
	"context"

	"coshkey_tree/internal/model"
	"coshkey_tree/internal/pkg/errors"
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
		return nil, internalerrors.ErrNotFound
	}

	return s.repository.GetUser(ctx, id)
}
