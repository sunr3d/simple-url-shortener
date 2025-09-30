package services

import "context"

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=Shortener --output=../../../mocks --filename=mock_shortener.go --with-expecter
type Shortener interface {
	ShortenLink(ctx context.Context, originalURL string) (string, string, error)
	FollowLink(ctx context.Context, code string) (string, error)
}
