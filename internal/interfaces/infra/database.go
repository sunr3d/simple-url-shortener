package infra

import (
	"context"

	"github.com/sunr3d/simple-url-shortener/models"
)

type Database interface {
	Create(ctx context.Context, link *models.Link) error
	GetLink(ctx context.Context, code string) (*models.Link, error)
}
