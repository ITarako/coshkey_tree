package user

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/ITarako/coshkey_tree/internal/database"
	"github.com/ITarako/coshkey_tree/internal/model"
	internalErrors "github.com/ITarako/coshkey_tree/internal/pkg/errors"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) GetUser(ctx context.Context, id int32) (*model.User, error) {
	sb := database.StatementBuilder.
		Select("id", "username").
		From("\"user\"").
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	user := new(model.User)
	err = r.db.QueryRowxContext(ctx, query, args...).StructScan(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrors.ErrNotFound
		}
		return nil, errors.Wrap(err, "db.QueryRowxContext()")
	}

	return user, nil
}
