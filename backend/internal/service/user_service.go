package service

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"thinh/gin-app/internal/dto"
	"thinh/gin-app/internal/model"
	"thinh/gin-app/internal/repository"
	apperrors "thinh/gin-app/pkg/errors"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
	jwtExpire int
}

func NewUserService(repo *repository.UserRepository, jwtSecret string, jwtExpire int) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret, jwtExpire: jwtExpire}
}

func (s *UserService) Create(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if existing != nil {
		return nil, apperrors.NewConflict("Email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}

	now := time.Now()
	user := &model.User{
		ID:           uuid.New().String(),
		FullName:     req.FullName,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		AvatarURL:    req.AvatarURL,
		AddressText:  req.AddressText,
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, apperrors.NewInternal(err)
	}

	return s.userToResponse(user), nil
}

func (s *UserService) GetByID(id string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if user == nil {
		return nil, apperrors.NewNotFound("User not found")
	}
	return s.userToResponse(user), nil
}

func (s *UserService) GetAll(page, pageSize int) ([]dto.UserResponse, int, error) {
	offset := (page - 1) * pageSize
	users, total, err := s.repo.FindAll(page, pageSize, offset)
	if err != nil {
		return nil, 0, apperrors.NewInternal(err)
	}

	resp := make([]dto.UserResponse, len(users))
	for i, u := range users {
		resp[i] = *s.userToResponse(&u)
	}
	return resp, total, nil
}

func (s *UserService) Update(id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if user == nil {
		return nil, apperrors.NewNotFound("User not found")
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Email != nil {
		// Check if new email already taken
		existing, err := s.repo.FindByEmail(*req.Email)
		if err != nil {
			return nil, apperrors.NewInternal(err)
		}
		if existing != nil && existing.ID != id {
			return nil, apperrors.NewConflict("Email already in use")
		}
		user.Email = *req.Email
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.AddressText != nil {
		user.AddressText = req.AddressText
	}

	if err := s.repo.Update(user); err != nil {
		return nil, apperrors.NewInternal(err)
	}

	return s.userToResponse(user), nil
}

func (s *UserService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return apperrors.NewInternal(err)
	}
	return nil
}

func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}
	if user == nil {
		return nil, apperrors.NewUnauthorized("Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.NewUnauthorized("Invalid email or password")
	}

	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return nil, apperrors.NewInternal(err)
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *s.userToResponse(user),
	}, nil
}

func (s *UserService) userToResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		Phone:       user.Phone,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		AddressText: user.AddressText,
		IsVerified:  user.IsVerified,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func (s *UserService) generateToken(userID, email string) (string, error) {
	// JWT token generation
	claims := &jwtClaims{
		UserID: userID,
		Email:  email,
	}
	return claims.Generate(s.jwtSecret, s.jwtExpire)
}
