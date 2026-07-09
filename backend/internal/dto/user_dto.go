package dto

import "time"

// --- Requests ---

type CreateUserRequest struct {
	FullName    string  `json:"full_name" binding:"required,min=1,max=255"`
	Phone       *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	Email       string  `json:"email" binding:"required,email,max=255"`
	Password    string  `json:"password" binding:"required,min=6,max=100"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	AddressText *string `json:"address_text,omitempty"`
}

type UpdateUserRequest struct {
	FullName    *string `json:"full_name,omitempty" binding:"omitempty,min=1,max=255"`
	Phone       *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	Email       *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	AddressText *string `json:"address_text,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// --- Responses ---

type UserResponse struct {
	ID          string    `json:"id"`
	FullName    string    `json:"full_name"`
	Phone       *string   `json:"phone,omitempty"`
	Email       string    `json:"email"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	AddressText *string   `json:"address_text,omitempty"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type PaginationParams struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

func (p *PaginationParams) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

func (p *PaginationParams) GetPageSize() int {
	if p.PageSize < 1 {
		return 20
	}
	return p.PageSize
}

func (p *PaginationParams) Offset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}
