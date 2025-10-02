package infra

import (
	"context"

	"github.com/sunr3d/simple-url-shortener/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=Database --output=../../../mocks --filename=mock_database.go --with-expecter
type Database interface {
	Create(ctx context.Context, link *models.Link) error
	GetLink(ctx context.Context, code string) (*models.Link, error)

	RecordClick(ctx context.Context, click models.ClickAnalytics) error
}
