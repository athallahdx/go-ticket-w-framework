package main

import (
	"go-ticket/internal/config"
	"go-ticket/internal/handler"
	"go-ticket/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SetupRouter(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	adminUserHandler *handler.AdminUserHandler,
	authHandler *handler.AuthHandler,
	cfg *config.Config,
) {
	router.Static("/uploads", "./uploads")

	api := router.Group("/api")

	// public
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/register", authHandler.Register)

	// authenticated
	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		auth.GET("/me", authHandler.GetProfile)
		auth.PUT("/me/password", authHandler.ChangePassword)
		auth.PUT("/me", userHandler.UpdateProfile)
	}

	// admin
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	admin.Use(middleware.RoleMiddleware("admin"))
	{
		admin.GET("/users", adminUserHandler.GetAllUsers)
		admin.GET("/users/:id", adminUserHandler.GetUserByID)
		admin.PUT("/users/:id", adminUserHandler.UpdateUser)
		admin.PATCH("/users/:id/role", adminUserHandler.UpdateRole)
		admin.DELETE("/users/:id", adminUserHandler.DeleteUser)
	}

	log.Info().Msg("✅ Routes configured successfully")
}
