package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"

	"github.com/sunr3d/simple-url-shortener/internal/config"
	"github.com/sunr3d/simple-url-shortener/internal/interfaces/infra"
	"github.com/sunr3d/simple-url-shortener/models"
)

const (
	queryCreate = `INSERT INTO urls(code, original_url) values ($1, $2)`
	queryRead   = `SELECT id, code, original_url, created_at FROM urls WHERE code = $1`
)

var _ infra.Database = (*postgresRepo)(nil)

type postgresRepo struct {
	db *dbpg.DB
}

func New(ctx context.Context, cfg config.DBConfig) (infra.Database, error) {
	db, err := dbpg.New(cfg.DSN, nil, &dbpg.Options{})
	if err != nil {
		return nil, fmt.Errorf("не удалось создать подключение к БД: %w", err)
	}

	pCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := db.Master.PingContext(pCtx); err != nil {
		_ = db.Master.Close()
		return nil, fmt.Errorf("таймаут пинг к БД: %w", err)
	}

	return &postgresRepo{db: db}, nil
}

func (r *postgresRepo) Create(ctx context.Context, link *models.Link) error {
	_, err := r.db.ExecWithRetry(
		ctx,
		retry.Strategy{Attempts: 3},
		queryCreate,
		link.Code,
		link.Original,
	)

	return err
}

func (r *postgresRepo) GetLink(ctx context.Context, code string) (*models.Link, error) {
	row := r.db.QueryRowContext(
		ctx,
		queryRead,
		code,
	)

	var l models.Link
	if err := row.Scan(
		&l.ID,
		&l.Code,
		&l.Original,
		&l.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("не удалось получить ссылку из БД: %w", err)
	}

	return &l, nil
}
