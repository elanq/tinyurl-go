package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/elanq/tinyurl-go/model"
	_ "github.com/lib/pq"
)

type URL interface {
	Create(context.Context, model.URL) (*model.URL, error)
	Update(context.Context, model.URL) error
	FindByShortUrl(ctx context.Context, url string) (*model.URL, error)
}

type url struct {
	db *sql.DB
}

func (u *url) Create(ctx context.Context, url model.URL) (*model.URL, error) {

	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO url (user_id, short_url, long_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", url.UserId, url.ShortURL, url.LongURL, time.Now(), time.Now())
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return nil, e
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &url, nil
}
func (u *url) Update(ctx context.Context, url model.URL) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "UPDATE url SET long_url = ? updated_at = now() WHERE user_id = ? AND short_url = ?", url.LongURL, url.UserId, url.ShortURL)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return e
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}
func (u *url) FindByShortUrl(ctx context.Context, url string) (*model.URL, error) {
	row := u.db.QueryRowContext(ctx, "SELECT * FROM url WHERE deleted_at is not null and short_url = $1", url)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var urlModel model.URL
	if err := row.Scan(&urlModel.ID, &urlModel.UserId, &urlModel.ShortURL, &urlModel.LongURL, &urlModel.CreatedAt, &urlModel.UpdatedAt, &urlModel.DeletedAt); err != nil {
		return nil, err
	}
	return &urlModel, nil
}

func NewURL(db *sql.DB) URL {
	return &url{
		db: db,
	}
}
