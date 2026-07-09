package service

import (
	"time"

	"github.com/google/uuid"

	"thinh/gin-app/internal/dto"
	"thinh/gin-app/internal/model"
	"thinh/gin-app/internal/repository"
	apperrors "thinh/gin-app/pkg/errors"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) Create(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	now := time.Now()
	post := &model.Post{
		ID:         uuid.New().String(),
		UserID:     userID,
		GroupID:    req.GroupID,
		Type:       req.Type,
		Title:      req.Title,
		Content:    req.Content,
		Location:   req.Location,
		ImageURLs:  req.ImageURLs,
		IsResolved: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, apperrors.NewInternal(err)
	}

	return s.postToResponse(post, ""), nil
}

func (s *PostService) GetByID(id string) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if post == nil {
		return nil, apperrors.NewNotFound("Post not found")
	}
	return s.postToResponse(post, ""), nil
}

func (s *PostService) GetAll(page, pageSize int, groupID, postType string) ([]dto.PostResponse, int, error) {
	offset := (page - 1) * pageSize
	posts, total, err := s.postRepo.FindAll(page, pageSize, offset, groupID, postType)
	if err != nil {
		return nil, 0, apperrors.NewInternal(err)
	}

	resp := make([]dto.PostResponse, len(posts))
	for i, p := range posts {
		resp[i] = *s.postToResponse(&p, "")
	}
	return resp, total, nil
}

func (s *PostService) Update(id string, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if post == nil {
		return nil, apperrors.NewNotFound("Post not found")
	}

	if req.Type != nil {
		post.Type = *req.Type
	}
	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Location != nil {
		post.Location = req.Location
	}
	if req.ImageURLs != nil {
		post.ImageURLs = req.ImageURLs
	}
	if req.IsResolved != nil {
		post.IsResolved = *req.IsResolved
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, apperrors.NewInternal(err)
	}

	return s.postToResponse(post, ""), nil
}

func (s *PostService) Delete(id string) error {
	if err := s.postRepo.Delete(id); err != nil {
		return apperrors.NewInternal(err)
	}
	return nil
}

func (s *PostService) MarkResolved(id string) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if post == nil {
		return nil, apperrors.NewNotFound("Post not found")
	}

	if err := s.postRepo.MarkResolved(id); err != nil {
		return nil, apperrors.NewInternal(err)
	}

	post.IsResolved = true
	return s.postToResponse(post, ""), nil
}

func (s *PostService) postToResponse(post *model.Post, authorName string) *dto.PostResponse {
	return &dto.PostResponse{
		ID:         post.ID,
		UserID:     post.UserID,
		GroupID:    post.GroupID,
		Type:       post.Type,
		Title:      post.Title,
		Content:    post.Content,
		Location:   post.Location,
		ImageURLs:  post.ImageURLs,
		IsResolved: post.IsResolved,
		ExpiresAt:  post.ExpiresAt,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
		AuthorName: authorName,
	}
}
