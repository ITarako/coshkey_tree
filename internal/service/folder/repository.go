package folder

import (
	"context"
	"database/sql"
	"fmt"

	"coshkey_tree/internal/database"
	"coshkey_tree/internal/model"
	"coshkey_tree/internal/pkg/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) RootFolder(ctx context.Context, userId int) (*model.Folder, error) {
	sb := database.StatementBuilder.
		Select("id", "id_user", "id_parent", "lft", "rgt", "depth", "title", "is_active", "is_project").
		From("folder").
		Where(sq.And{
			sq.Eq{"id_user": userId},
			sq.Eq{"depth": 1},
			sq.Gt{"id_parent": 0},
		}).
		OrderBy("lft", "rgt").
		Limit(1)

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	folder := new(model.Folder)
	err = r.db.QueryRowxContext(ctx, query, args...).StructScan(folder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalerrors.ErrNotFound
		}
		return nil, errors.Wrap(err, "db.QueryRowxContext()")
	}

	return folder, nil
}

func (r Repository) Children(ctx context.Context, parentLft, parentRgt int) ([]model.Folder, error) {
	sb := database.StatementBuilder.
		Select("id", "id_user", "id_parent", "lft", "rgt", "depth", "title", "is_active", "is_project").
		From("folder").
		Where(sq.And{
			sq.Gt{"lft": parentLft},
			sq.Lt{"rgt": parentRgt},
			sq.Eq{"is_active": true},
		}).OrderBy("title")

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	var res []model.Folder
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "db.SelectContext()")
	}

	return res, nil
}

func (r Repository) AdminFolders(ctx context.Context) ([]model.Folder, error) {
	sb := database.StatementBuilder.
		Select("id", "id_user", "id_parent", "lft", "rgt", "depth", "title", "is_active", "is_project").
		From("folder").
		Where(sq.Gt{"depth": 0}).
		OrderBy("title")

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	var res []model.Folder
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "db.SelectContext()")
	}

	return res, nil
}

func (r Repository) UserFoldersByRight(ctx context.Context, userId, right int) ([]model.Folder, error) {
	query := `
		SELECT
			id, id_user, id_parent, lft, rgt, depth, title, is_active, is_project
		FROM folder
		WHERE depth > 0
			AND is_active = TRUE
			AND (id_user = $1 OR
				id IN (SELECT model_id FROM right_control WHERE id_user = $2 AND right_control.right = $3 AND model_class = $4))
		ORDER BY title
`

	var res []model.Folder
	err := r.db.SelectContext(ctx, &res, query, userId, userId, right, model.FolderClassName)
	if err != nil {
		return nil, errors.Wrap(err, "db.SelectContext()")
	}

	return res, nil
}

func (r Repository) UserShareFolder(ctx context.Context, userId int) ([]model.Folder, error) {
	condition, err := r.getWhereFilter(ctx, userId)
	if err != nil || len(condition) == 0 {
		return nil, err
	}

	text := `
		SELECT
			id, id_user, id_parent, lft, rgt, depth, title, is_active, is_project
		FROM folder
		WHERE depth > 0
			AND is_active = TRUE
			AND (%s)
		ORDER BY title
`
	query := fmt.Sprintf(text, condition)

	var res []model.Folder
	err = r.db.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, errors.Wrap(err, "db.SelectContext()")
	}

	return res, nil
}

func (r Repository) getWhereFilter(ctx context.Context, userId int) (string, error) {
	query := `
		SELECT
			lft, rgt, id_user
		FROM folder
			LEFT JOIN key ON folder.id = key.id_folder
		WHERE folder.is_active = TRUE
			AND (folder.id IN (SELECT model_id FROM right_control WHERE id_user = $1 AND model_class = $2)
				OR folder.id IN (SELECT id_folder
				          FROM key
				          WHERE id IN (SELECT model_id FROM right_control WHERE id_user = $3 AND model_class = $4))
				)
`

	rows, err := r.db.QueryContext(ctx, query, userId, model.FolderClassName, userId, model.KeyClassName)
	if err != nil {
		return "", errors.Wrap(err, "db.QueryContext()")
	}

	defer func() {
		_ = rows.Close()
	}()

	var condition string
	for rows.Next() {
		var lft, rgt, idUser int
		if err = rows.Scan(&lft, &rgt, &idUser); err == nil {
			text := "(lft <= %d AND rgt >= %d AND id_user = %d) OR "
			condition += fmt.Sprintf(text, lft, rgt, idUser)
		}
	}

	if len(condition) > 0 {
		condition = condition[:len(condition)-4]
	}

	return condition, nil
}

func (r Repository) UserFavorite(ctx context.Context, userId int) (map[int]model.Folder, error) {
	query := `
		SELECT
			folder.id as id,
			folder.id_user as id_user,
			id_parent,
			lft,
			rgt,
			depth,
			folder.title as title,
			folder.is_active as is_active,
			is_project,
			CASE WHEN ff.id IS NULL THEN false ELSE true END as is_favorite,
			count(fk.id) as count_favorite_keys
		FROM folder
			LEFT JOIN favorites as ff on ff.id_item = folder.id and ff.type = $1 and ff.is_active = true and ff.id_user = $2
			LEFT JOIN key on key.id_folder = folder.id and key.is_active = true and (key.login IS NOT NULL OR key.password IS NOT NULL)
			LEFT JOIN favorites as fk on fk.id_item = key.id and fk.type = $3 and fk.is_active = true and fk.id_user = $4
		WHERE depth > 0
		  AND folder.is_active = true
		GROUP BY folder.id, ff.id
`

	rows, err := r.db.QueryxContext(ctx, query, model.FavoriteTypeFolder, userId, model.FavoriteTypeKey, userId)
	if err != nil {
		return nil, errors.Wrap(err, "db.QueryxContext()")
	}

	defer func() {
		_ = rows.Close()
	}()

	var res = make(map[int]model.Folder)
	for rows.Next() {
		var favorite model.Folder
		if err = rows.StructScan(&favorite); err == nil {
			res[favorite.Id] = favorite
		}
	}

	return res, nil
}
