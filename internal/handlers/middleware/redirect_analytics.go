package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/simple-url-shortener/internal/interfaces/services"
	"github.com/sunr3d/simple-url-shortener/models"
)

func RedirectAnalytics(s services.Shortener) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		click := models.ClickAnalytics{
			UserAgent:  c.Request.UserAgent(),
			IP:         c.ClientIP(),
			Referrer:   c.Request.Referer(),
			OccurredAt: time.Now(),
		}

		c.Next()

		if c.Writer.Status() == http.StatusFound {
			click.Code = c.Param("short_url")

			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)

			go func() {
				defer cancel()
				if err := s.RecordClick(ctx, click); err != nil {
					zlog.Logger.Warn().Err(err).Msg("не удалось записать клик")
				}
			}()
		}
	}
}
