package rbac

import (
	"context"

	"github.com/pkg/errors"

	"coshkey_tree/internal/model"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) GetRolesByUser(ctx context.Context, userId int) (map[string]model.AuthItem, error) {
	roles, err := s.repository.GetRolesByUser(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "repository.GetRolesByUser()")
	}

	res := make(map[string]model.AuthItem)
	for _, role := range roles {
		res[role.Name] = role
	}

	return res, nil
}

func (s Service) CheckRole(ctx context.Context, roleName string, userId int) (bool, error) {
	roles, err := s.GetRolesByUser(ctx, userId)
	if err != nil {
		return false, errors.Wrap(err, "service.GetRolesByUser()")
	}

	_, ok := roles[roleName]

	return ok, nil
}
