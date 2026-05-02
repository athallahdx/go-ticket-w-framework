package main

import (
	"go-ticket/internal/config"
	"go-ticket/internal/handler"
	"go-ticket/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SetupRouter(router *gin.Engine, userHandler *handler.UserHandler, authHandler *handler.AuthHandler, cfg *config.Config) {
	router.Static("/uploads", "./uploads")

	api := router.Group("/api")

	api.POST("/login", authHandler.Login)
	api.POST("/register", authHandler.Register)

	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		auth.GET("/profile", authHandler.GetProfile)
		auth.PUT("/profile/change-password", authHandler.ChangePassword)
		auth.PUT("/profile/update", userHandler.UpdateProfile)
	}

	log.Info().Msg("API Routes configured successfully")
}
