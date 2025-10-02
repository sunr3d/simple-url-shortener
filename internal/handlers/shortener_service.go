package httphandlers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/simple-url-shortener/models"
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
		zlog.Logger.Error().Err(err).Msg("svc.ShortenLink")
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
		if strings.Contains(err.Error(), "не найдена") {
			zlog.Logger.Error().Err(err).Msg("svc.FollowLink")
			c.JSON(http.StatusNotFound, ginext.H{"error": "ссылка не найдена"})
			return
		}
		zlog.Logger.Error().Err(err).Msg("svc.FollowLink")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "внутерняя ошибка сервера"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

func (h *Handler) getAnalytics(c *ginext.Context) {
	code := strings.TrimSpace(c.Param("short_url"))
	if code == "" || len(code) > 32 {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный short_url"})
		return
	}

	var tr models.TimeRange
	if s := strings.TrimSpace(c.Query("from")); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "from: YYYY-MM-DD"})
			return
		}
		tr.From = t
	}

	if s := strings.TrimSpace(c.Query("to")); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "to: YYYY-MM-DD"})
			return
		}
		tr.To = t
	}

	group := strings.TrimSpace(c.Query("group"))

	res, err := h.svc.GetAnalytics(c.Request.Context(), code, tr, group)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			zlog.Logger.Error().Err(err).Msg("svc.GetAnalytics")
			c.JSON(http.StatusNotFound, ginext.H{"error": "ссылка не найдена"})
			return
		}
		zlog.Logger.Error().Err(err).Msg("svc.GetAnalytics")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "внутренняя ошибка сервера"})
		return
	}

	out := analyticsResp{Total: res.Total}
	if len(res.ByDay) > 0 {
		out.ByDay = make([]analyticsByDay, 0, len(res.ByDay))
		for _, b := range res.ByDay {
			out.ByDay = append(out.ByDay, analyticsByDay{Date: b.Date, Count: b.Count})
		}
	}

	if len(res.ByMonth) > 0 {
		out.ByMonth = make([]analyticsByMonth, 0, len(res.ByMonth))
		for _, b := range res.ByMonth {
			out.ByMonth = append(out.ByMonth, analyticsByMonth{Year: b.Year, Month: b.Month, Count: b.Count})
		}
	}

	if len(res.ByUA) > 0 {
		out.ByUA = make([]analyticsByUA, 0, len(res.ByUA))
		for _, b := range res.ByUA {
			out.ByUA = append(out.ByUA, analyticsByUA{UA: b.UserAgent, Count: b.Count})
		}
	}

	c.JSON(http.StatusOK, out)
}
