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

	queryRecordClick = `INSERT INTO url_clicks(url_id, occurred_at, ip_address, user_agent, referrer) SELECT u.id, $5, $3, $2, $4 FROM urls u WHERE u.code = $1`

	queryClicksTotal = `
	SELECT COUNT(*)
	FROM url_clicks uc
	JOIN urls u ON u.id=uc.url_id
	WHERE u.code = $1
		AND ($2::timestamp IS NULL OR uc.occurred_at >= $2)
		AND ($3::timestamp IS NULL OR uc.occurred_at <= $3)`
	queryClicksByDay = `
	SELECT to_char(date_trunc('day', uc.occurred_at), 'YYYY-MM-DD') AS bucket, COUNT(*)
	FROM url_clicks uc
	JOIN urls u ON u.id=uc.url_id
	WHERE u.code = $1
		AND ($2::timestamp IS NULL OR uc.occurred_at >= $2)
		AND ($3::timestamp IS NULL OR uc.occurred_at <= $3)
	GROUP BY bucket
	ORDER BY bucket`
	queryClicksByMonth = `
	SELECT EXTRACT(YEAR FROM uc.occurred_at)::int AS y,
		EXTRACT(MONTH FROM uc.occurred_at)::int AS m,
		COUNT(*)
	FROM url_clicks uc
	JOIN urls u ON u.id=uc.url_id
	WHERE u.code = $1
		AND ($2::timestamp IS NULL OR uc.occurred_at >= $2)
		AND ($3::timestamp IS NULL OR uc.occurred_at <= $3)
	GROUP BY y,m
	ORDER BY y,m`
	queryClicksByUA = `
	SELECT COALESCE(NULLIF(uc.user_agent,''), 'unknown') AS ua, COUNT(*)
	FROM url_clicks uc
	JOIN urls u ON u.id=uc.url_id
	WHERE u.code = $1
		AND ($2::timestamp IS NULL OR uc.occurred_at >= $2)
		AND ($3::timestamp IS NULL OR uc.occurred_at <= $3)
	GROUP BY ua
	ORDER BY COUNT(*) DESC`
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

func (r *postgresRepo) RecordClick(ctx context.Context, click models.ClickAnalytics) error {
	_, err := r.db.ExecWithRetry(
		ctx,
		retry.Strategy{Attempts: 3},
		queryRecordClick,
		click.Code,
		click.UserAgent,
		click.IP,
		click.Referrer,
		click.OccurredAt,
	)

	return err
}

func (r *postgresRepo) GetTotal(ctx context.Context, code string, tr models.TimeRange) (int64, error) {
	fromPtr, toPtr := r.trBounds(tr)

	row := r.db.QueryRowContext(ctx, queryClicksTotal, code, fromPtr, toPtr)

	var out int64
	if err := row.Scan(&out); err != nil {
		return 0, fmt.Errorf("не удалось получить total аналитику: %w", err)
	}

	return out, nil
}

func (r *postgresRepo) GetByDay(ctx context.Context, code string, tr models.TimeRange) ([]models.ClicksByDay, error) {
	fromPtr, toPtr := r.trBounds(tr)

	rows, err := r.db.QueryContext(ctx, queryClicksByDay, code, fromPtr, toPtr)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить аналитику по дням: %w", err)
	}
	defer rows.Close()

	var out []models.ClicksByDay

	for rows.Next() {
		var date string
		var count int64
		if err := rows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании БД: %w", err)
		}
		out = append(out, models.ClicksByDay{Date: date, Count: count})
	}

	return out, rows.Err()
}

func (r *postgresRepo) GetByMonth(ctx context.Context, code string, tr models.TimeRange) ([]models.ClicksByMonth, error) {
	fromPtr, toPtr := r.trBounds(tr)

	rows, err := r.db.QueryContext(ctx, queryClicksByMonth, code, fromPtr, toPtr)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить аналитику по месяцам: %w", err)
	}
	defer rows.Close()

	var out []models.ClicksByMonth

	for rows.Next() {
		var year, month int
		var count int64
		if err := rows.Scan(&year, &month, &count); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании БД: %w", err)
		}
		out = append(out, models.ClicksByMonth{Year: year, Month: month, Count: count})
	}

	return out, rows.Err()
}

func (r *postgresRepo) GetByUserAgent(ctx context.Context, code string, tr models.TimeRange) ([]models.ClicksByUserAgent, error) {
	fromPtr, toPtr := r.trBounds(tr)

	rows, err := r.db.QueryContext(ctx, queryClicksByUA, code, fromPtr, toPtr)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить аналитику по UserAgent: %w", err)
	}
	defer rows.Close()

	var out []models.ClicksByUserAgent

	for rows.Next() {
		var ua string
		var count int64
		if err := rows.Scan(&ua, &count); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании БД: %w", err)
		}
		out = append(out, models.ClicksByUserAgent{UserAgent: ua, Count: count})
	}

	return out, rows.Err()
}

// trBounds - хелпер для установки границ временного периода.
func (r *postgresRepo) trBounds(tr models.TimeRange) (fromPtr, toPtr any) {
	if !tr.From.IsZero() {
		fromPtr = tr.From
	} else {
		fromPtr = nil
	}

	if !tr.To.IsZero() {
		toPtr = tr.To
	} else {
		toPtr = nil
	}

	return
}
