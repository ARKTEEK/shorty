package models

import (
	"time"
)

type Link struct {
	ID          int32     `json:"id"`
	OriginalUrl string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	UserID      int32     `json:"user_id"`
	Visits      int32     `json:"visits"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

type CreateLinkRequest struct {
	OriginalUrl string `json:"original_url"`
	UserId      int32  `json:"user_id"`
}
