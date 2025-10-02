package models

import "time"

type ClickAnalytics struct {
	ID         int64
	Code       string
	UserAgent  string
	IP         string
	Referrer   string
	OccurredAt time.Time
}
