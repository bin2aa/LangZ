package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thinh/gin-app/internal/dto"
	"thinh/gin-app/internal/service"
	"thinh/gin-app/pkg/response"
)

type GroupController struct {
	service *service.GroupService
}

func NewGroupController(svc *service.GroupService) *GroupController {
	return &GroupController{service: svc}
}

// Create creates a new group
// @Summary Create a group
// @Description Create a new group (requires authentication)
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateGroupRequest true "Group details"
// @Success 201 {object} response.APIResponse{data=dto.GroupResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /groups [post]
func (c *GroupController) Create(ctx *gin.Context) {
	var req dto.CreateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	userID, _ := ctx.Get("user_id")
	group, err := c.service.Create(userID.(string), &req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Created(ctx, group)
}

// GetByID returns a group by ID
// @Summary Get group by ID
// @Description Get a single group by its UUID
// @Tags Groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} response.APIResponse{data=dto.GroupResponse}
// @Failure 404 {object} response.APIResponse
// @Router /groups/{id} [get]
func (c *GroupController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	group, err := c.service.GetByID(id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, group)
}

// GetAll returns a paginated list of groups
// @Summary List all groups
// @Description Get a paginated list of all groups
// @Tags Groups
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.APIResponse
// @Router /groups [get]
func (c *GroupController) GetAll(ctx *gin.Context) {
	var pagination dto.PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	groups, total, err := c.service.GetAll(pagination.GetPage(), pagination.GetPageSize())
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithPagination(ctx, groups, pagination.GetPage(), pagination.GetPageSize(), total)
}

// Update updates a group's information
// @Summary Update group
// @Description Update group details (requires authentication)
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Param request body dto.UpdateGroupRequest true "Updated group details"
// @Success 200 {object} response.APIResponse{data=dto.GroupResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /groups/{id} [put]
func (c *GroupController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.UpdateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	group, err := c.service.Update(id, &req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, group)
}

// Delete removes a group
// @Summary Delete group
// @Description Delete a group (requires authentication)
// @Tags Groups
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Success 204 "No Content"
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /groups/{id} [delete]
func (c *GroupController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.NoContent(ctx)
}

// Join allows a user to join a group
// @Summary Join a group
// @Description Join a group (requires authentication)
// @Tags Groups
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /groups/{id}/join [post]
func (c *GroupController) Join(ctx *gin.Context) {
	groupID := ctx.Param("id")
	userID, _ := ctx.Get("user_id")

	if err := c.service.JoinGroup(groupID, userID.(string)); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "Successfully joined the group"})
}

// Leave allows a user to leave a group
// @Summary Leave a group
// @Description Leave a group (requires authentication)
// @Tags Groups
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /groups/{id}/leave [post]
func (c *GroupController) Leave(ctx *gin.Context) {
	groupID := ctx.Param("id")
	userID, _ := ctx.Get("user_id")

	if err := c.service.LeaveGroup(groupID, userID.(string)); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "Successfully left the group"})
}

// GetMembers returns all members of a group
// @Summary Get group members
// @Description Get all members of a specific group
// @Tags Groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /groups/{id}/members [get]
func (c *GroupController) GetMembers(ctx *gin.Context) {
	groupID := ctx.Param("id")
	members, err := c.service.GetMembers(groupID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, members)
}
