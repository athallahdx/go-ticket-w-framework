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
	adminOrganizerHandler *handler.AdminOrganizerHandler,
	adminEventHandler *handler.AdminEventHandler,
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
		// User Management
		users := admin.Group("/users")
		users.GET("/", adminUserHandler.GetAllUsers)
		users.GET("/:id", adminUserHandler.GetUserByID)
		users.PUT("/:id", adminUserHandler.UpdateUser)
		users.PATCH("/:id/role", adminUserHandler.UpdateRole)
		users.DELETE("/:id", adminUserHandler.DeleteUser)

		// Organizer Management
		organizers := admin.Group("/organizers")
		organizers.POST("/", adminOrganizerHandler.CreateOrganizer)
		organizers.GET("/", adminOrganizerHandler.GetAllOrganizers)
		organizers.GET("/:id", adminOrganizerHandler.GetOrganizerByID)
		organizers.PUT("/:id", adminOrganizerHandler.UpdateOrganizer)
		organizers.DELETE("/:id", adminOrganizerHandler.DeleteOrganizer)

		// Event Management
		events := admin.Group("/events")
		events.POST("/", adminEventHandler.CreateEvent)
		events.GET("/", adminEventHandler.GetAllEvents)
		events.GET("/:id", adminEventHandler.GetEventByID)
		events.PUT("/:id", adminEventHandler.UpdateEvent)
		events.DELETE("/:id", adminEventHandler.DeleteEvent)
	}

	log.Info().Msg("✅ Routes configured successfully")
}
