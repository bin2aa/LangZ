package repository

import (
	"database/sql"
	"fmt"
	"time"

	"thinh/gin-app/internal/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *model.Post) error {
	query := `
		INSERT INTO posts (id, user_id, group_id, type, title, content, location, image_urls, is_resolved, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, ST_GeogFromText($7), $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query,
		post.ID, post.UserID, post.GroupID, post.Type, post.Title, post.Content,
		post.Location, post.ImageURLs, post.IsResolved, post.ExpiresAt,
		post.CreatedAt, post.UpdatedAt,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (r *PostRepository) FindByID(id string) (*model.Post, error) {
	query := `SELECT p.id, p.user_id, p.group_id, p.type, p.title, p.content, 
	           ST_AsText(p.location) as location, p.image_urls, p.is_resolved, p.expires_at, p.created_at, p.updated_at
	           FROM posts p WHERE p.id = $1`
	post := &model.Post{}
	var groupID, location sql.NullString
	var expiresAt sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&post.ID, &post.UserID, &groupID, &post.Type, &post.Title, &post.Content,
		&location, &post.ImageURLs, &post.IsResolved, &expiresAt, &post.CreatedAt, &post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find post by id: %w", err)
	}
	if groupID.Valid {
		post.GroupID = &groupID.String
	}
	if location.Valid {
		post.Location = &location.String
	}
	if expiresAt.Valid {
		post.ExpiresAt = &expiresAt.Time
	}
	return post, nil
}

func (r *PostRepository) FindAll(page, pageSize, offset int, groupID, postType string) ([]model.Post, int, error) {
	// Build count query
	countQuery := `SELECT COUNT(*) FROM posts p WHERE 1=1`
	query := `SELECT p.id, p.user_id, p.group_id, p.type, p.title, p.content, 
	           ST_AsText(p.location) as location, p.image_urls, p.is_resolved, p.expires_at, p.created_at, p.updated_at,
	           u.full_name as author_name
	           FROM posts p JOIN users u ON p.user_id = u.id WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if groupID != "" {
		countQuery += fmt.Sprintf(" AND p.group_id = $%d", argIdx)
		query += fmt.Sprintf(" AND p.group_id = $%d", argIdx)
		args = append(args, groupID)
		argIdx++
	}
	if postType != "" {
		countQuery += fmt.Sprintf(" AND p.type = $%d", argIdx)
		query += fmt.Sprintf(" AND p.type = $%d", argIdx)
		args = append(args, postType)
		argIdx++
	}

	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count posts: %w", err)
	}

	query += fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		var groupID, location sql.NullString
		var expiresAt sql.NullTime
		var authorName string
		if err := rows.Scan(
			&p.ID, &p.UserID, &groupID, &p.Type, &p.Title, &p.Content,
			&location, &p.ImageURLs, &p.IsResolved, &expiresAt, &p.CreatedAt, &p.UpdatedAt,
			&authorName,
		); err != nil {
			return nil, 0, fmt.Errorf("scan post: %w", err)
		}
		if groupID.Valid {
			p.GroupID = &groupID.String
		}
		if location.Valid {
			p.Location = &location.String
		}
		if expiresAt.Valid {
			p.ExpiresAt = &expiresAt.Time
		}
		posts = append(posts, p)
	}
	return posts, total, nil
}

func (r *PostRepository) Update(post *model.Post) error {
	query := `UPDATE posts SET type = $1, title = $2, content = $3, location = ST_GeogFromText($4), 
	           image_urls = $5, is_resolved = $6, updated_at = $7 WHERE id = $8`
	_, err := r.db.Exec(query, post.Type, post.Title, post.Content, post.Location,
		post.ImageURLs, post.IsResolved, time.Now(), post.ID)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}
	return nil
}

func (r *PostRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostRepository) MarkResolved(id string) error {
	_, err := r.db.Exec(`UPDATE posts SET is_resolved = true, updated_at = $1 WHERE id = $2`, time.Now(), id)
	return err
}
