package httphandlers

import (
	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	svc *shortenersvc.Service
}

func New(svc *shortenersvc.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterHandlers() *ginext.Engine {
	router := ginext.New()
	router.Use(ginext.Logger(), ginext.Recovery())

	router.POST("/shorten", h.shortenLink)
	router.GET("/s/:short_url", h.redirect)
	router.GET("/analytics/:short_url", h.getAnalytics)

	return router
}