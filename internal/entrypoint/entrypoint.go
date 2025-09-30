package entrypoint

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sunr3d/simple-url-shortener/internal/config"
)

func Run(cfg *config.Config) error {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Инфраслой

	// Сервисный слой

	// REST API (HTTP) + Middleware

	// Server

	<-appCtx.Done()
	return nil
}
