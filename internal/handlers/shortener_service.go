package httphandlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) shortenLink(c *ginext.Context) {
	var req shortenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный JSON"})
		return
	}

	raw := strings.TrimSpace(req.URL)
	if raw == "" {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "необходимо указать URL для сокращения"})
		return
	}
	if len(raw) > 2048 {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "слишком длинный URL"})
		return
	}

	if _, err := url.ParseRequestURI(raw); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "неверный синтаксис URL"})
		return
	}

	code, shortURL, err := h.svc.ShortenLink(c.Request.Context(), raw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, shortenResp{
		Code:     code,
		ShortURL: shortURL,
	})
}

func (h *Handler) redirect(c *ginext.Context) {
	code := strings.TrimSpace(c.Param("short_url"))
	if code == "" || len(code) > 32 {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный short_url"})
		return
	}

	originalURL, err := h.svc.FollowLink(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "ссылка не найдена"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

func (h *Handler) getAnalytics(c *ginext.Context) {}
