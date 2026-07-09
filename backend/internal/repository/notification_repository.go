package repository

import (
	"database/sql"
	"fmt"

	"thinh/gin-app/internal/model"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *model.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, post_id, message, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`
	return r.db.QueryRow(query,
		notification.ID, notification.UserID, notification.PostID, notification.Message,
		notification.IsRead, notification.CreatedAt,
	).Scan(&notification.ID, &notification.CreatedAt)
}

func (r *NotificationRepository) FindByUserID(userID string, page, pageSize, offset int) ([]model.Notification, int, error) {
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count notifications: %w", err)
	}

	query := `SELECT id, user_id, post_id, message, is_read, created_at 
	           FROM notifications WHERE user_id = $1 
	           ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		var postID sql.NullString
		if err := rows.Scan(&n.ID, &n.UserID, &postID, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan notification: %w", err)
		}
		if postID.Valid {
			n.PostID = &postID.String
		}
		notifications = append(notifications, n)
	}
	return notifications, total, nil
}

func (r *NotificationRepository) MarkAsRead(id string) error {
	_, err := r.db.Exec(`UPDATE notifications SET is_read = true WHERE id = $1`, id)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(userID string) error {
	_, err := r.db.Exec(`UPDATE notifications SET is_read = true WHERE user_id = $1`, userID)
	return err
}
