package httphandlers

import (
	"github.com/wb-go/wbf/ginext"

	"github.com/sunr3d/simple-url-shortener/internal/handlers/middleware"
	"github.com/sunr3d/simple-url-shortener/internal/interfaces/services"
)

type Handler struct {
	svc services.Shortener
}

func New(svc services.Shortener) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterHandlers() *ginext.Engine {
	router := ginext.New()
	router.Use(ginext.Logger(), ginext.Recovery())

	router.POST("/shorten", h.shortenLink)
	router.GET("/s/:short_url", middleware.RedirectAnalytics(h.svc), h.redirect)
	router.GET("/analytics/:short_url", h.getAnalytics)

	return router
}
