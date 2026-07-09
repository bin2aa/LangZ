package dto

import "time"

type CreatePostRequest struct {
	GroupID   *string  `json:"group_id,omitempty"`
	Type      string   `json:"type" binding:"required,oneof=lost_found give_away alert general service_request"`
	Title     string   `json:"title" binding:"required,min=1,max=255"`
	Content   string   `json:"content" binding:"required,min=1"`
	Location  *string  `json:"location,omitempty"`
	ImageURLs []string `json:"image_urls,omitempty"`
}

type UpdatePostRequest struct {
	Type       *string  `json:"type,omitempty" binding:"omitempty,oneof=lost_found give_away alert general service_request"`
	Title      *string  `json:"title,omitempty" binding:"omitempty,min=1,max=255"`
	Content    *string  `json:"content,omitempty" binding:"omitempty,min=1"`
	Location   *string  `json:"location,omitempty"`
	ImageURLs  []string `json:"image_urls,omitempty"`
	IsResolved *bool    `json:"is_resolved,omitempty"`
}

type PostResponse struct {
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
	AuthorName string     `json:"author_name,omitempty"`
}
