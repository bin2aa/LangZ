package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thinh/gin-app/internal/controller"
	"thinh/gin-app/internal/middleware"
	"thinh/gin-app/internal/repository"
	"thinh/gin-app/internal/service"
	"thinh/gin-app/pkg/logger"
)

// scalarHTML is the Scalar API Reference UI HTML page.
// It loads the generated OpenAPI spec from /swagger.json and renders the interactive UI.
const scalarHTML = `<!doctype html>
<html>
<head>
    <title>Neighborhood API - Scalar</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
    <script id="api-reference" data-url="/swagger.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

func Setup(
	db *gin.Engine,
	userRepo *repository.UserRepository,
	groupRepo *repository.GroupRepository,
	postRepo *repository.PostRepository,
	notificationRepo *repository.NotificationRepository,
	jwtSecret string,
	jwtExpire int,
) {
	// Initialize services
	userService := service.NewUserService(userRepo, jwtSecret, jwtExpire)
	groupService := service.NewGroupService(groupRepo, userRepo)
	postService := service.NewPostService(postRepo)
	notificationService := service.NewNotificationService(notificationRepo)

	// Initialize controllers
	userCtrl := controller.NewUserController(userService)
	groupCtrl := controller.NewGroupController(groupService)
	postCtrl := controller.NewPostController(postService)
	notifCtrl := controller.NewNotificationController(notificationService)

	// Middleware
	db.Use(middleware.CORSMiddleware())
	db.Use(gin.Logger())
	db.Use(gin.Recovery())

	// API docs
	db.StaticFile("/swagger.json", "docs/swagger.json")
	db.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, scalarHTML)
	})

	api := db.Group("/api/v1")
	{
		// Public routes
		api.POST("/auth/login", userCtrl.Login)
		api.POST("/users", userCtrl.Register)

		// Public: list and view resources
		api.GET("/users", userCtrl.GetAll)
		api.GET("/users/:id", userCtrl.GetByID)
		api.GET("/groups", groupCtrl.GetAll)
		api.GET("/groups/:id", groupCtrl.GetByID)
		api.GET("/groups/:id/members", groupCtrl.GetMembers)
		api.GET("/posts", postCtrl.GetAll)
		api.GET("/posts/:id", postCtrl.GetByID)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthRequired(jwtSecret))
		{
			// User profile
			protected.GET("/profile", userCtrl.GetProfile)
			protected.PUT("/users/:id", userCtrl.Update)
			protected.DELETE("/users/:id", userCtrl.Delete)

			// Groups
			protected.POST("/groups", groupCtrl.Create)
			protected.PUT("/groups/:id", groupCtrl.Update)
			protected.DELETE("/groups/:id", groupCtrl.Delete)
			protected.POST("/groups/:id/join", groupCtrl.Join)
			protected.POST("/groups/:id/leave", groupCtrl.Leave)

			// Posts
			protected.POST("/posts", postCtrl.Create)
			protected.PUT("/posts/:id", postCtrl.Update)
			protected.DELETE("/posts/:id", postCtrl.Delete)
			protected.PATCH("/posts/:id/resolve", postCtrl.MarkResolved)

			// Notifications
			protected.GET("/notifications", notifCtrl.GetAll)
			protected.PATCH("/notifications/:id/read", notifCtrl.MarkAsRead)
			protected.PATCH("/notifications/read-all", notifCtrl.MarkAllAsRead)
		}
	}

	logger.Info("Routes initialized successfully")
}
