package model

import (
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID        int       `json:"id"`
	UserID	uuid.UUID `json:"user_id"`
	LongURL   string    `json:"long_url"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

