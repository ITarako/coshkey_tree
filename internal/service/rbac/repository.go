package rbac

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/ITarako/coshkey_tree/internal/database"
	"github.com/ITarako/coshkey_tree/internal/model/rbac"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) GetRolesByUser(ctx context.Context, userId int32) ([]model.AuthItem, error) {
	sb := database.StatementBuilder.
		Select("ai.name", "ai.type", "ai.description").
		From("auth_item as ai").
		LeftJoin("auth_assignment as aa on aa.item_name=ai.name").
		Where(sq.And{
			sq.Eq{"aa.user_id": userId},
			sq.Eq{"ai.type": model.TYPE_ROLE},
		})

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	var res []model.AuthItem
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "db.SelectContext()")
	}

	return res, nil
}
