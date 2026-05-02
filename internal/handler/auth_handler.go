package handler

import (
	"go-ticket/internal/config"
	"go-ticket/internal/domain"
	"go-ticket/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	authService domain.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService domain.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to register user")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().Str("email", req.Email).Msg("User registered successfully")

	c.JSON(http.StatusCreated, h.toUserResponse(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Login failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Info().Str("email", req.Email).Int64("user_id", user.ID).Msg("User logged in successfully")

	c.JSON(http.StatusOK, dto.AuthResponse{
		Token: token,
		User:  h.toUserResponse(user),
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		log.Error().Err(err).Int64("user_id", userID).Msg("Failed to get user profile")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.toUserResponse(user))
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		log.Error().Err(err).Int64("user_id", userID).Msg("Failed to change password")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("user_id", userID).Msg("User changed password successfully")

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

func getUserID(c *gin.Context) (int64, error) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, gin.Error{Err: nil, Type: gin.ErrorTypePublic}
	}

	return val.(int64), nil
}

func (h *AuthHandler) buildProfileURL(profile string) string {
	if profile == "" {
		return ""
	}
	return h.cfg.BaseURL + "/uploads/" + profile
}

func (h *AuthHandler) toUserResponse(user *domain.User) dto.UserResponse {
	if user == nil {
		return dto.UserResponse{}
	}

	return dto.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Role:    user.Role,
		Profile: h.buildProfileURL(user.Profile),
	}
}
