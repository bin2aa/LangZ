package dto

import "time"

type CreateGroupRequest struct {
	Name           string  `json:"name" binding:"required,min=1,max=255"`
	Description    *string `json:"description,omitempty"`
	CenterLocation *string `json:"center_location,omitempty"`
	RadiusMeters   int     `json:"radius_meters" binding:"required,min=1"`
}

type UpdateGroupRequest struct {
	Name           *string `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Description    *string `json:"description,omitempty"`
	CenterLocation *string `json:"center_location,omitempty"`
	RadiusMeters   *int    `json:"radius_meters,omitempty" binding:"omitempty,min=1"`
}

type GroupResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    *string   `json:"description,omitempty"`
	CenterLocation *string   `json:"center_location,omitempty"`
	RadiusMeters   int       `json:"radius_meters"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	MemberCount    int       `json:"member_count,omitempty"`
}

type GroupMemberResponse struct {
	UserID   string    `json:"user_id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	JoinedAt time.Time `json:"joined_at"`
}
