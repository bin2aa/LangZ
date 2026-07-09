package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thinh/gin-app/internal/dto"
	"thinh/gin-app/internal/service"
	"thinh/gin-app/pkg/response"
)

type PostController struct {
	service *service.PostService
}

func NewPostController(svc *service.PostService) *PostController {
	return &PostController{service: svc}
}

// Create creates a new post
// @Summary Create a post
// @Description Create a new post (requires authentication)
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePostRequest true "Post details"
// @Success 201 {object} response.APIResponse{data=dto.PostResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /posts [post]
func (c *PostController) Create(ctx *gin.Context) {
	var req dto.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	userID, _ := ctx.Get("user_id")
	post, err := c.service.Create(userID.(string), &req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Created(ctx, post)
}

// GetByID returns a post by ID
// @Summary Get post by ID
// @Description Get a single post by its UUID
// @Tags Posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.APIResponse{data=dto.PostResponse}
// @Failure 404 {object} response.APIResponse
// @Router /posts/{id} [get]
func (c *PostController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	post, err := c.service.GetByID(id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, post)
}

// GetAll returns a paginated list of posts
// @Summary List all posts
// @Description Get a paginated list of posts, optionally filtered by group_id and type
// @Tags Posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param group_id query string false "Filter by group ID"
// @Param type query string false "Filter by post type"
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.APIResponse
// @Router /posts [get]
func (c *PostController) GetAll(ctx *gin.Context) {
	var pagination dto.PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	groupID := ctx.Query("group_id")
	postType := ctx.Query("type")

	posts, total, err := c.service.GetAll(pagination.GetPage(), pagination.GetPageSize(), groupID, postType)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithPagination(ctx, posts, pagination.GetPage(), pagination.GetPageSize(), total)
}

// Update updates a post
// @Summary Update post
// @Description Update a post (requires authentication)
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body dto.UpdatePostRequest true "Updated post details"
// @Success 200 {object} response.APIResponse{data=dto.PostResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /posts/{id} [put]
func (c *PostController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	post, err := c.service.Update(id, &req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, post)
}

// Delete removes a post
// @Summary Delete post
// @Description Delete a post (requires authentication)
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 204 "No Content"
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /posts/{id} [delete]
func (c *PostController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.NoContent(ctx)
}

// MarkResolved marks a post as resolved
// @Summary Mark post as resolved
// @Description Mark a post as resolved (requires authentication)
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} response.APIResponse{data=dto.PostResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /posts/{id}/resolve [patch]
func (c *PostController) MarkResolved(ctx *gin.Context) {
	id := ctx.Param("id")
	post, err := c.service.MarkResolved(id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, post)
}
