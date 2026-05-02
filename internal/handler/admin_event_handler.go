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

type AdminEventHandler struct {
	adminEventService domain.AdminEventService
	cfg               *config.Config
}

func NewAdminEventHandler(adminEventService domain.AdminEventService, cfg *config.Config) *AdminEventHandler {
	return &AdminEventHandler{
		adminEventService: adminEventService,
		cfg:               cfg,
	}
}

func (h *AdminEventHandler) CreateEvent(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &domain.Event{
		OrganizerID: req.OrganizerID,
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		Thumbnail:   req.Thumbnail,
		Date:        req.Date,
	}

	if err := h.adminEventService.CreateEvent(event); err != nil {
		log.Error().Err(err).Msg("Failed to create event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("event_id", event.ID).Msg("Event created successfully by admin")

	c.JSON(http.StatusCreated, h.toEventResponse(event))
}

func (h *AdminEventHandler) GetAllEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	events, total, err := h.adminEventService.GetAllEvents(page, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all events")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var eventResponses []dto.EventResponse
	for _, event := range events {
		eventResponses = append(eventResponses, h.toEventResponse(event))
	}

	c.JSON(http.StatusOK, dto.EventListResponse{
		Events: eventResponses,
		Total:  int64(total),
	})
}

func (h *AdminEventHandler) GetEventByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	event, err := h.adminEventService.GetEventByID(id)
	if err != nil {
		log.Error().Err(err).Int64("event_id", id).Msg("Failed to get event by ID")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.toEventResponse(event))
}

func (h *AdminEventHandler) UpdateEvent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &domain.Event{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		Thumbnail:   req.Thumbnail,
		Date:        req.Date,
	}

	if err := h.adminEventService.UpdateEvent(event); err != nil {
		log.Error().Err(err).Int64("event_id", id).Msg("Failed to update event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("event_id", id).Msg("Event updated successfully by admin")

	updatedEvent, _ := h.adminEventService.GetEventByID(id)
	c.JSON(http.StatusOK, h.toEventResponse(updatedEvent))
}

func (h *AdminEventHandler) DeleteEvent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	if err := h.adminEventService.DeleteEvent(id); err != nil {
		log.Error().Err(err).Int64("event_id", id).Msg("Failed to delete event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("event_id", id).Msg("Event deleted successfully by admin")

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func (h *AdminEventHandler) toEventResponse(event *domain.Event) dto.EventResponse {
	if event == nil {
		return dto.EventResponse{}
	}

	return dto.EventResponse{
		ID:          event.ID,
		OrganizerID: event.OrganizerID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		Thumbnail:   event.Thumbnail, // Add a builder method later if needed
		Date:        event.Date,
		CreatedAt:   event.CreatedAt,
	}
}
