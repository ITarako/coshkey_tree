package tree

import (
	"github.com/jmoiron/sqlx"

	"github.com/ITarako/coshkey_tree/internal/model"
	"github.com/ITarako/coshkey_tree/internal/service/folder"
	"github.com/ITarako/coshkey_tree/internal/service/rbac"
	"github.com/ITarako/coshkey_tree/internal/service/user"
)

type Result struct {
	Tree     string `json:"tree"`
	Favorite string `json:"favorite"`
}

type Service struct {
	db            *sqlx.DB
	UserService   user.Service
	FolderService folder.Service
	RbacService   rbac.Service
}

func NewService(db *sqlx.DB) Service {
	return Service{
		UserService:   user.NewService(user.NewRepository(db)),
		FolderService: folder.NewService(folder.NewRepository(db)),
		RbacService:   rbac.NewService(rbac.NewRepository(db)),
	}
}

func (s *Service) Generate(user *model.User, folderId int32, ownTree bool) Result {
	res := Result{}

	res.Tree = s.generateMain()
	res.Favorite = s.generateFavorite()

	return res
}

func (s *Service) generateMain() string {
	return ""
}

func (s *Service) generateFavorite() string {
	return ""
}
