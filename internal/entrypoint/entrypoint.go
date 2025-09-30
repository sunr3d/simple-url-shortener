package entrypoint

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wb-go/wbf/zlog"
	
	"github.com/sunr3d/simple-url-shortener/internal/config"
	"github.com/sunr3d/simple-url-shortener/internal/infra/postgres"
)

func Run(cfg *config.Config) error {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Инфраслой
	repo, err := postgres.New(appCtx, cfg.DB)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("posgres.New")
		return fmt.Errorf("postgres.New(): %w", err)
	}
	// Сервисный слой

	// REST API (HTTP) + Middleware

	// Server

	<-appCtx.Done()
	return nil
}
