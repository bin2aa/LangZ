package model

import "time"

type User struct {
	ID           string    `json:"id"`
	FullName     string    `json:"full_name"`
	Phone        *string   `json:"phone,omitempty"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	AddressText  *string   `json:"address_text,omitempty"`
	Location     *string   `json:"location,omitempty"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
