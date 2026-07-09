package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thinh/gin-app/internal/dto"
	"thinh/gin-app/internal/service"
	"thinh/gin-app/pkg/response"
)

type NotificationController struct {
	service *service.NotificationService
}

func NewNotificationController(svc *service.NotificationService) *NotificationController {
	return &NotificationController{service: svc}
}

// GetAll returns notifications for the current user
// @Summary Get my notifications
// @Description Get a paginated list of notifications for the authenticated user
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.PaginatedResponse
// @Failure 401 {object} response.APIResponse
// @Router /notifications [get]
func (c *NotificationController) GetAll(ctx *gin.Context) {
	var pagination dto.PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	userID, _ := ctx.Get("user_id")
	notifications, total, err := c.service.GetByUserID(userID.(string), pagination.GetPage(), pagination.GetPageSize())
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithPagination(ctx, notifications, pagination.GetPage(), pagination.GetPageSize(), total)
}

// MarkAsRead marks a notification as read
// @Summary Mark notification as read
// @Description Mark a single notification as read (requires authentication)
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /notifications/{id}/read [patch]
func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.MarkAsRead(id); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// MarkAllAsRead marks all notifications as read for the current user
// @Summary Mark all notifications as read
// @Description Mark all notifications as read for the authenticated user
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /notifications/read-all [patch]
func (c *NotificationController) MarkAllAsRead(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	if err := c.service.MarkAllAsRead(userID.(string)); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "All notifications marked as read"})
}
