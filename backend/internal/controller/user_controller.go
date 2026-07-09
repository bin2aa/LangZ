package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thinh/gin-app/internal/dto"
	"thinh/gin-app/internal/service"
	"thinh/gin-app/pkg/response"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(svc *service.UserService) *UserController {
	return &UserController{service: svc}
}

// Register creates a new user account
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User registration details"
// @Success 201 {object} response.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} response.APIResponse
// @Router /users [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := c.service.Create(&req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Created(ctx, user)
}

// Login authenticates a user
// @Summary Login
// @Description Authenticate user with email and password, returns JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.APIResponse{data=dto.LoginResponse}
// @Failure 400 {object} response.APIResponse
// @Router /auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := c.service.Login(&req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, resp)
}

// GetByID returns a user by ID
// @Summary Get user by ID
// @Description Get a single user by their UUID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.APIResponse{data=dto.UserResponse}
// @Failure 404 {object} response.APIResponse
// @Router /users/{id} [get]
func (c *UserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.service.GetByID(id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, user)
}

// GetAll returns a paginated list of users
// @Summary List all users
// @Description Get a paginated list of all users
// @Tags Users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.APIResponse
// @Router /users [get]
func (c *UserController) GetAll(ctx *gin.Context) {
	var pagination dto.PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	users, total, err := c.service.GetAll(pagination.GetPage(), pagination.GetPageSize())
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithPagination(ctx, users, pagination.GetPage(), pagination.GetPageSize(), total)
}

// Update updates a user's information
// @Summary Update user
// @Description Update user details (requires authentication)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "Updated user details"
// @Success 200 {object} response.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /users/{id} [put]
func (c *UserController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := c.service.Update(id, &req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, user)
}

// Delete removes a user account
// @Summary Delete user
// @Description Delete a user account (requires authentication)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /users/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.NoContent(ctx)
}

// GetProfile returns the current authenticated user's profile
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=dto.UserResponse}
// @Failure 401 {object} response.APIResponse
// @Router /profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	user, err := c.service.GetByID(userID.(string))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, user)
}
