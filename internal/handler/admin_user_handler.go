package handler

import (
	"go-ticket/internal/config"
	"go-ticket/internal/domain"
	"go-ticket/internal/dto"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AdminUserHandler struct {
	adminUserService domain.AdminUserService
	cfg              *config.Config
}

func NewAdminUserHandler(adminUserService domain.AdminUserService, cfg *config.Config) *AdminUserHandler {
	return &AdminUserHandler{
		adminUserService: adminUserService,
		cfg:              cfg,
	}
}

func (h *AdminUserHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, err := h.adminUserService.GetAllUsers(page, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, h.toUserResponse(user))
	}

	c.JSON(http.StatusOK, dto.UserListResponse{
		Users: userResponses,
		Total: int64(total),
	})
}

func (h *AdminUserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := h.adminUserService.GetUserByID(id)
	if err != nil {
		log.Error().Err(err).Int64("user_id", id).Msg("Failed to get user by ID")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.toUserResponse(user))
}

func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		ID:      id,
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Profile: req.Profile,
		Role:    req.Role,
	}

	if err := h.adminUserService.UpdateUser(user); err != nil {
		log.Error().Err(err).Int64("user_id", id).Msg("Failed to update user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("user_id", id).Msg("User updated successfully by admin")

	// Fetch the updated user
	updatedUser, _ := h.adminUserService.GetUserByID(id)
	c.JSON(http.StatusOK, h.toUserResponse(updatedUser))
}

func (h *AdminUserHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adminUserService.UpdateRole(id, req.Role); err != nil {
		log.Error().Err(err).Int64("user_id", id).Str("role", req.Role).Msg("Failed to update user role")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("user_id", id).Str("role", req.Role).Msg("User role updated successfully by admin")

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := h.adminUserService.DeleteUser(id); err != nil {
		log.Error().Err(err).Int64("user_id", id).Msg("Failed to delete user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("user_id", id).Msg("User deleted successfully by admin")

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *AdminUserHandler) buildProfileURL(profile string) string {
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

func (h *AdminUserHandler) toUserResponse(user *domain.User) dto.UserResponse {
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
