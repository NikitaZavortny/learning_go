package main

import (
	"auth-server/internal/config"
	"auth-server/internal/handler"
	"auth-server/internal/middleware"
	"auth-server/internal/repository"
	"auth-server/internal/service"
	"auth-server/pkg/database"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresDB(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repository
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	mailService := service.NewMailService()
	authService := service.NewAuthService(userRepo, mailService, cfg.JWTSecret)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, cfg.JWTSecret)

	// Initialize router
	router := gin.Default()

	router.Use(middleware.CORS(&middleware.CORSConfig{
		AllowOrigins:     cfg.GetCORSAllowOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/activate/:link", authHandler.Activate)
	}
	// Protected routes
	protectedGroup := router.Group("/api")
	protectedGroup.Use(middleware.AuthMiddleware(userRepo, cfg.JWTSecret))
	{
		protectedGroup.GET("/profile", authHandler.Profile)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
