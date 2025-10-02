package services

import (
	"context"

	"github.com/sunr3d/simple-url-shortener/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=Shortener --output=../../../mocks --filename=mock_shortener.go --with-expecter
type Shortener interface {
	ShortenLink(ctx context.Context, originalURL string) (string, string, error)
	FollowLink(ctx context.Context, code string) (string, error)

	RecordClick(ctx context.Context, click models.ClickAnalytics) error
	GetAnalytics(ctx context.Context, code string, tr models.TimeRange, group string) (models.Analytics, error)
}
