package models

import "time"

type Link struct {
	ID        int64
	Code      string
	Original  string
	CreatedAt time.Time
}
