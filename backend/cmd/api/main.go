package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"thinh/gin-app/config"
	"thinh/gin-app/internal/repository"
	"thinh/gin-app/internal/routes"
	"thinh/gin-app/pkg/database"
	"thinh/gin-app/pkg/logger"
)

// @title           Neighborhood API
// @version         1.0
// @description     API for neighborhood social network
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization

func main() {
	// Initialize logger
	logger.Init()

	// Parse config path flag
	configPath := flag.String("config", "config/config.dev.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Override config path from env if set
	if envConfig := os.Getenv("CONFIG_PATH"); envConfig != "" {
		cfg, err = config.Load(envConfig)
		if err != nil {
			log.Fatalf("Failed to load config from env: %v", err)
		}
	}

	logger.Info("Connecting to database...")
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Info("Database connected successfully")

	// Run migrations
	logger.Info("Running database migrations...")
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	logger.Info("Migrations completed")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	postRepo := repository.NewPostRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// Setup Gin router
	ginMode := gin.DebugMode
	if os.Getenv("GIN_MODE") == "release" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	router := gin.New()

	// Setup routes
	routes.Setup(
		router,
		userRepo,
		groupRepo,
		postRepo,
		notificationRepo,
		cfg.JWT.Secret,
		cfg.JWT.ExpireHour,
	)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Info("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
