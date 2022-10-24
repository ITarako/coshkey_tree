package folder

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

func (r Repository) GetFolder(ctx context.Context, id int32) (*model.Folder, error) {
	sb := database.StatementBuilder.
		Select("id", "id_user", "id_parent", "lft", "rgt", "depth", "title", "is_active", "is_project").
		From("folder").
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	folder := new(model.Folder)
	err = r.db.QueryRowxContext(ctx, query, args...).StructScan(folder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrors.ErrNotFound
		}
		return nil, errors.Wrap(err, "db.QueryRowxContext()")
	}

	return folder, nil
}
