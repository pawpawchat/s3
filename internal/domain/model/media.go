package model

import "time"

type Media struct {
	ID          int64
	OwnerID     int64
	Description string
	ContentType string
	FileExt     string
	CreatedAt   time.Time
}
