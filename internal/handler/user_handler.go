package handler

import (
	"net/http"
	"strings"

	"go-ticket/internal/config"
	"go-ticket/internal/domain"
	"go-ticket/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	userService domain.UserService
	cfg         *config.Config
}

func NewUserHandler(userService domain.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		cfg:         cfg,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetProfileByID(userID)
	if err != nil {
		log.Error().Err(err).Int64("user_id", userID).Msg("Failed to get user profile")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.toUserResponse(user))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var input domain.UpdateProfileInput

	contentType := c.GetHeader("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
			return
		}

		if name := c.PostForm("name"); name != "" {
			input.Name = &name
		}
		if phone := c.PostForm("phone"); phone != "" {
			input.Phone = &phone
		}

		fileHeader, err := c.FormFile("profile")
		if err != nil && err.Error() != "http: no such file" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
			return
		}

		user, err := h.userService.UpdateProfile(userID, input, fileHeader)
		if err != nil {
			log.Error().Err(err).Int64("user_id", userID).Msg("Failed to update profile with file")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Info().Int64("user_id", userID).Msg("User profile updated successfully with file")

		c.JSON(http.StatusOK, h.toUserResponse(user))
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		input.Name = &req.Name
	}
	if req.Phone != "" {
		input.Phone = &req.Phone
	}

	user, err := h.userService.UpdateProfile(userID, input, nil)
	if err != nil {
		log.Error().Err(err).Int64("user_id", userID).Msg("Failed to update profile")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("user_id", userID).Msg("User profile updated successfully")

	c.JSON(http.StatusOK, h.toUserResponse(user))
}

func (h *UserHandler) buildProfileURL(profile string) string {
	if profile == "" {
		return ""
	}
	if strings.HasPrefix(profile, "http") {
		return profile
	}
	if strings.HasPrefix(profile, "uploads/") {
		return h.cfg.BaseURL + "/" + profile
	}
	return h.cfg.BaseURL + "/uploads/" + profile
}

func (h *UserHandler) toUserResponse(user *domain.User) dto.UserResponse {
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
