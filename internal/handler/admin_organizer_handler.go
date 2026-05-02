package handler

import (
	"go-ticket/internal/config"
	"go-ticket/internal/domain"
	"go-ticket/internal/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AdminOrganizerHandler struct {
	adminOrganizerService domain.AdminOrganizerService
	cfg                   *config.Config
}

func NewAdminOrganizerHandler(adminOrganizerService domain.AdminOrganizerService, cfg *config.Config) *AdminOrganizerHandler {
	return &AdminOrganizerHandler{
		adminOrganizerService: adminOrganizerService,
		cfg:                   cfg,
	}
}

func (h *AdminOrganizerHandler) CreateOrganizer(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	userID, _ := strconv.ParseInt(c.PostForm("user_id"), 10, 64)

	organizer := &domain.Organizer{
		UserID:      userID,
		CompanyName: c.PostForm("company_name"),
		Phone:       c.PostForm("phone"),
		Email:       c.PostForm("email"),
		Description: c.PostForm("description"),
		City:        c.PostForm("city"),
		Province:    c.PostForm("province"),
		IsVerified:  false, // default state
	}

	fileHeader, _ := c.FormFile("logo")

	if err := h.adminOrganizerService.CreateOrganizer(organizer, fileHeader); err != nil {
		log.Error().Err(err).Msg("Failed to create organizer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("organizer_id", organizer.ID).Msg("Organizer created successfully by admin")

	c.JSON(http.StatusCreated, h.toOrganizerResponse(organizer))
}

func (h *AdminOrganizerHandler) GetAllOrganizers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	organizers, total, err := h.adminOrganizerService.GetAllOrganizers(page, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all organizers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var organizerResponses []dto.OrganizerResponse
	for _, org := range organizers {
		organizerResponses = append(organizerResponses, h.toOrganizerResponse(org))
	}

	c.JSON(http.StatusOK, dto.OrganizerListResponse{
		Organizers: organizerResponses,
		Total:      int64(total),
	})
}

func (h *AdminOrganizerHandler) GetOrganizerByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organizer ID format"})
		return
	}

	organizer, err := h.adminOrganizerService.GetOrganizerByID(id)
	if err != nil {
		log.Error().Err(err).Int64("organizer_id", id).Msg("Failed to get organizer by ID")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.toOrganizerResponse(organizer))
}

func (h *AdminOrganizerHandler) UpdateOrganizer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organizer ID format"})
		return
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	organizer := &domain.Organizer{
		ID:          id,
		CompanyName: c.PostForm("company_name"),
		Phone:       c.PostForm("phone"),
		Email:       c.PostForm("email"),
		Description: c.PostForm("description"),
		City:        c.PostForm("city"),
		Province:    c.PostForm("province"),
	}

	fileHeader, _ := c.FormFile("logo")

	if err := h.adminOrganizerService.UpdateOrganizer(organizer, fileHeader); err != nil {
		log.Error().Err(err).Int64("organizer_id", id).Msg("Failed to update organizer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("organizer_id", id).Msg("Organizer updated successfully by admin")

	updatedOrganizer, _ := h.adminOrganizerService.GetOrganizerByID(id)
	c.JSON(http.StatusOK, h.toOrganizerResponse(updatedOrganizer))
}

func (h *AdminOrganizerHandler) DeleteOrganizer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organizer ID format"})
		return
	}

	if err := h.adminOrganizerService.DeleteOrganizer(id); err != nil {
		log.Error().Err(err).Int64("organizer_id", id).Msg("Failed to delete organizer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("organizer_id", id).Msg("Organizer deleted successfully by admin")

	c.JSON(http.StatusOK, gin.H{"message": "Organizer deleted successfully"})
}

func (h *AdminOrganizerHandler) toOrganizerResponse(org *domain.Organizer) dto.OrganizerResponse {
	if org == nil {
		return dto.OrganizerResponse{}
	}

	return dto.OrganizerResponse{
		ID:          org.ID,
		UserID:      org.UserID,
		CompanyName: org.CompanyName,
		Phone:       org.Phone,
		Email:       org.Email,
		Logo:        org.Logo, // You could add buildLogoURL here if similar to buildProfileURL
		Description: org.Description,
		City:        org.City,
		Province:    org.Province,
		IsVerified:  org.IsVerified,
		VerifiedAt:  org.VerifiedAt,
		CreatedAt:   org.CreatedAt,
	}
}
