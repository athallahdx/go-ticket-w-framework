package main

import (
	"fmt"

	"go-ticket/internal/config"
	"go-ticket/internal/handler"
	"go-ticket/internal/repository"
	"go-ticket/internal/service"
	"go-ticket/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	logFile, err := logger.Init("storage/logs/app.log")
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logFile.Close()

	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Error open database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting sql.DB")
	}
	defer sqlDB.Close()

	userRepo := repository.NewUserRepository(db)
	organizerRepo := repository.NewOrganizerRepository(db)
	eventRepo := repository.NewEventRepository(db)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	adminUserService := service.NewAdminUserService(userRepo)
	adminOrganizerService := service.NewAdminOrganizerService(organizerRepo)
	adminEventService := service.NewAdminEventService(eventRepo)

	userHandler := handler.NewUserHandler(userService, cfg)
	authHandler := handler.NewAuthHandler(authService, cfg)
	adminUserHandler := handler.NewAdminUserHandler(adminUserService, cfg)
	adminOrganizerHandler := handler.NewAdminOrganizerHandler(adminOrganizerService, cfg)
	adminEventHandler := handler.NewAdminEventHandler(adminEventService, cfg)

	log.Info().Msg("✅ MySQL connected Successfully!")
	log.Info().Str("port", cfg.Port).Msg("✅ Server running")

	router := gin.Default()
	SetupRouter(
		router,
		userHandler,
		adminUserHandler,
		adminOrganizerHandler,
		adminEventHandler,
		authHandler,
		cfg,
	)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
