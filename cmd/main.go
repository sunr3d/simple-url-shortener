package main

import (
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/simple-url-shortener/internal/entrypoint"
	"github.com/sunr3d/simple-url-shortener/internal/config"
)

func main() {
	zlog.Init()
	zlog.Logger.Info().Msg("запуск приложения...")

	zlog.Logger.Info().Msg("загрузка конфига...")
	cfg, err := config.GetConfig(".env")
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("config.GetConfig")
	}
	zlog.Logger.Info().Msg("конфиг успешно загружен...")

	if err := entrypoint.Run(cfg); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("entrypoint.Run")
	}
}
