package shortenersvc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/sunr3d/simple-url-shortener/internal/interfaces/infra"
	"github.com/sunr3d/simple-url-shortener/internal/interfaces/services"
	"github.com/sunr3d/simple-url-shortener/models"
)

const codeLen = 5

var _ services.Shortener = (*shortenerService)(nil)

type shortenerService struct {
	repo    infra.Database
	baseURL string
}

func New(repo infra.Database, baseURL string) *shortenerService {
	return &shortenerService{
		repo:    repo,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *shortenerService) ShortenLink(ctx context.Context, originalURL string) (string, string, error) {
	u, err := url.Parse(originalURL)
	if err != nil {
		return "", "", fmt.Errorf("неверный URL")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", "", fmt.Errorf("неподдержваемый протокол URL")
	}

	code, err := generateCode(codeLen)
	if err != nil {
		return "", "", fmt.Errorf("не удалось сгенерировать код: %w", err)
	}

	if err := s.repo.Create(ctx, &models.Link{
		Code:     code,
		Original: originalURL,
	}); err != nil {
		return "", "", fmt.Errorf("repo.Create(): %w", err)
	}

	return code, s.baseURL + "/s/" + code, nil
}

func (s *shortenerService) FollowLink(ctx context.Context, code string) (string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", fmt.Errorf("code не может быть пустой")
	}

	link, err := s.repo.GetLink(ctx, code)
	if err != nil {
		return "", fmt.Errorf("repo.GetLink(): %w", err)
	}
	if link == nil {
		return "", fmt.Errorf("ссылка не найдена")
	}

	return link.Original, nil
}

func (s *shortenerService) RecordClick(ctx context.Context, click models.ClickAnalytics) error {
	if strings.TrimSpace(click.Code) == "" {
		return fmt.Errorf("click.Code не может быть пустым")
	}

	return s.repo.RecordClick(ctx, click)
}

// generateCode - генерирует URL совместимый рандомный код.
func generateCode(n int) (string, error) {
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	code := base64.RawURLEncoding.EncodeToString(b)
	if len(code) < n {
		return code, nil
	}

	return code[:n], nil
}
