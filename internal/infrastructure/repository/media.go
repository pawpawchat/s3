package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pawpawchat/s3/internal/domain/model"
)

type MediaRepository struct {
	db *sqlx.DB
}

func NewMediaRepository(db *sqlx.DB) *MediaRepository {
	return &MediaRepository{db}
}

func (r *MediaRepository) CreateMedia(ctx context.Context, media *model.Media) error {
	sql, args := squirrel.Insert("media").
		Columns("owner_id", "content_type", "file_ext").
		Values(media.OwnerID, media.ContentType, media.FileExt).
		PlaceholderFormat(squirrel.Dollar).
		MustSql()

	res, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	media.ID, err = res.LastInsertId()
	return err
}
