package model

import "time"

type Post struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	GroupID    *string    `json:"group_id,omitempty"`
	Type       string     `json:"type"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Location   *string    `json:"location,omitempty"`
	ImageURLs  []string   `json:"image_urls,omitempty"`
	IsResolved bool       `json:"is_resolved"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
