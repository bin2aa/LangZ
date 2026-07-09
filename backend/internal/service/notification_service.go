package service

import (
	"github.com/google/uuid"

	"thinh/gin-app/internal/dto"
	"thinh/gin-app/internal/model"
	"thinh/gin-app/internal/repository"
	apperrors "thinh/gin-app/pkg/errors"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetByUserID(userID string, page, pageSize int) ([]dto.NotificationResponse, int, error) {
	offset := (page - 1) * pageSize
	notifications, total, err := s.repo.FindByUserID(userID, page, pageSize, offset)
	if err != nil {
		return nil, 0, apperrors.NewInternal(err)
	}

	resp := make([]dto.NotificationResponse, len(notifications))
	for i, n := range notifications {
		resp[i] = *s.notificationToResponse(&n)
	}
	return resp, total, nil
}

func (s *NotificationService) MarkAsRead(id string) error {
	return s.repo.MarkAsRead(id)
}

func (s *NotificationService) MarkAllAsRead(userID string) error {
	return s.repo.MarkAllAsRead(userID)
}

func (s *NotificationService) CreateNotification(userID, postID, message string) (*dto.NotificationResponse, error) {
	notification := &model.Notification{
		ID:      uuid.New().String(),
		UserID:  userID,
		PostID:  &postID,
		Message: message,
		IsRead:  false,
	}

	if err := s.repo.Create(notification); err != nil {
		return nil, apperrors.NewInternal(err)
	}

	return s.notificationToResponse(notification), nil
}

func (s *NotificationService) notificationToResponse(n *model.Notification) *dto.NotificationResponse {
	return &dto.NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		PostID:    n.PostID,
		Message:   n.Message,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}
}
