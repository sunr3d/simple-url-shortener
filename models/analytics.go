package models

import "time"

type TimeRange struct {
	From time.Time
	To   time.Time
}

type ClickAnalytics struct {
	ID         int64
	Code       string
	UserAgent  string
	IP         string
	Referrer   string
	OccurredAt time.Time
}

type ClicksByDay struct {
	Date  string
	Count int64
}

type ClicksByMonth struct {
	Year  int
	Month int
	Count int64
}

type ClicksByUserAgent struct {
	UserAgent string
	Count     int64
}

type Analytics struct {
	Total   int64
	ByDay   []ClicksByDay
	ByMonth []ClicksByMonth
	ByUA    []ClicksByUserAgent
}
