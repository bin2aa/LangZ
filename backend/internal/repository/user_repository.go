package repository

import (
	"database/sql"
	"fmt"
	"time"

	"thinh/gin-app/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (id, full_name, phone, email, password_hash, avatar_url, address_text, location, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query,
		user.ID, user.FullName, user.Phone, user.Email, user.PasswordHash,
		user.AvatarURL, user.AddressText, user.Location, user.IsVerified,
		user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	query := `SELECT id, full_name, phone, email, password_hash, avatar_url, address_text, 
	           ST_AsText(location) as location, is_verified, created_at, updated_at 
	           FROM users WHERE id = $1`
	user := &model.User{}
	var phone, avatarURL, addressText, location sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.FullName, &phone, &user.Email, &user.PasswordHash,
		&avatarURL, &addressText, &location, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if phone.Valid {
		user.Phone = &phone.String
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if addressText.Valid {
		user.AddressText = &addressText.String
	}
	if location.Valid {
		user.Location = &location.String
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	query := `SELECT id, full_name, phone, email, password_hash, avatar_url, address_text, 
	           ST_AsText(location) as location, is_verified, created_at, updated_at 
	           FROM users WHERE email = $1`
	user := &model.User{}
	var phone, avatarURL, addressText, location sql.NullString
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.FullName, &phone, &user.Email, &user.PasswordHash,
		&avatarURL, &addressText, &location, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if phone.Valid {
		user.Phone = &phone.String
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if addressText.Valid {
		user.AddressText = &addressText.String
	}
	if location.Valid {
		user.Location = &location.String
	}
	return user, nil
}

func (r *UserRepository) FindAll(page, pageSize, offset int) ([]model.User, int, error) {
	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// Fetch page
	query := `SELECT id, full_name, phone, email, password_hash, avatar_url, address_text, 
	           ST_AsText(location) as location, is_verified, created_at, updated_at 
	           FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var phone, avatarURL, addressText, location sql.NullString
		if err := rows.Scan(
			&u.ID, &u.FullName, &phone, &u.Email, &u.PasswordHash,
			&avatarURL, &addressText, &location, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		if phone.Valid {
			u.Phone = &phone.String
		}
		if avatarURL.Valid {
			u.AvatarURL = &avatarURL.String
		}
		if addressText.Valid {
			u.AddressText = &addressText.String
		}
		if location.Valid {
			u.Location = &location.String
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserRepository) Update(user *model.User) error {
	query := `UPDATE users SET full_name = $1, phone = $2, email = $3, avatar_url = $4, 
	           address_text = $5, updated_at = $6 WHERE id = $7`
	_, err := r.db.Exec(query, user.FullName, user.Phone, user.Email, user.AvatarURL,
		user.AddressText, time.Now(), user.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
