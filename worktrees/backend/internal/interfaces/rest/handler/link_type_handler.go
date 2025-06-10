package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openfoundry/oms/internal/domain/entity"
	"github.com/openfoundry/oms/internal/domain/repository"
	"github.com/openfoundry/oms/internal/domain/service"
	"github.com/openfoundry/oms/internal/interfaces/rest/middleware"
	"github.com/openfoundry/oms/internal/pkg/validator"
	"go.uber.org/zap"
)

// LinkTypeHandler handles link type related requests
type LinkTypeHandler struct {
	service *service.LinkTypeService
	logger  *zap.Logger
}

// NewLinkTypeHandler creates a new link type handler
func NewLinkTypeHandler(service *service.LinkTypeService, logger *zap.Logger) *LinkTypeHandler {
	return &LinkTypeHandler{
		service: service,
		logger:  logger,
	}
}

// List handles GET /api/v1/link-types
func (h *LinkTypeHandler) List(c *gin.Context) {
	// Parse query parameters
	filter := repository.LinkTypeFilter{
		PageSize: 20, // Default page size
	}

	// Parse source object type filter
	if sourceID := c.Query("source_object_type_id"); sourceID != "" {
		id, err := uuid.Parse(sourceID)
		if err == nil {
			filter.SourceObjectTypeID = &id
		}
	}

	// Parse target object type filter
	if targetID := c.Query("target_object_type_id"); targetID != "" {
		id, err := uuid.Parse(targetID)
		if err == nil {
			filter.TargetObjectTypeID = &id
		}
	}

	// Parse cardinality filter
	if cardinality := c.Query("cardinality"); cardinality != "" {
		card := entity.Cardinality(cardinality)
		if card.IsValid() {
			filter.Cardinality = &card
		}
	}

	// Parse pagination
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			if validatedSize, err := validator.ValidatePageSize(pageSize); err == nil {
				filter.PageSize = validatedSize
			}
		}
	}

	if cursor := c.Query("cursor"); cursor != "" {
		filter.PageCursor = cursor
	}

	// Parse sort
	if sortBy := c.Query("sort_by"); sortBy != "" {
		allowedFields := []string{"name", "created_at", "updated_at"}
		if field, err := validator.ValidateSortBy(sortBy, allowedFields); err == nil {
			filter.SortBy = field
		}
	}

	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		if order, err := validator.ValidateSortOrder(sortOrder); err == nil {
			filter.SortOrder = order
		}
	}

	// Get link types
	linkTypes, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to list link types", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve link types",
		})
		return
	}

	// Generate next cursor if needed
	var nextCursor string
	if len(linkTypes) == filter.PageSize {
		lastItem := linkTypes[len(linkTypes)-1]
		nextCursor = encodeCursor(lastItem.CreatedAt, lastItem.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": linkTypes,
		"pagination": gin.H{
			"next_cursor": nextCursor,
			"page_size":   filter.PageSize,
		},
	})
}

// Create handles POST /api/v1/link-types
func (h *LinkTypeHandler) Create(c *gin.Context) {
	var input service.CreateLinkTypeInput

	// Bind and validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Additional validation
	if err := validator.ValidateObjectTypeName(input.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid link type name",
			"details": err.Error(),
		})
		return
	}

	// Validate cardinality
	if !input.Cardinality.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cardinality value",
		})
		return
	}

	// Sanitize input to prevent XSS
	input.Name = validator.SanitizeString(input.Name)
	input.DisplayName = validator.SanitizeString(input.DisplayName)
	if input.Description != nil {
		sanitized := validator.SanitizeString(*input.Description)
		input.Description = &sanitized
	}
	if input.InverseDisplayName != nil {
		sanitized := validator.SanitizeString(*input.InverseDisplayName)
		input.InverseDisplayName = &sanitized
	}

	// Get user ID from context
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Create link type
	linkType, err := h.service.CreateLinkType(c.Request.Context(), input, userID)
	if err != nil {
		h.logger.Error("Failed to create link type",
			zap.String("user_id", userID),
			zap.String("name", input.Name),
			zap.Error(err))

		// Handle specific errors
		switch err {
		case entity.ErrLinkTypeNameExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": "Link type name already exists",
			})
		case entity.ErrCircularReference:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Link type would create a circular reference",
			})
		case entity.ErrObjectTypeNotFound:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Source or target object type not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create link type",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, linkType)
}

// Get handles GET /api/v1/link-types/:id
func (h *LinkTypeHandler) Get(c *gin.Context) {
	// Parse ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid link type ID",
		})
		return
	}

	// Get link type
	linkType, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == entity.ErrLinkTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Link type not found",
			})
			return
		}

		h.logger.Error("Failed to get link type",
			zap.String("id", id.String()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve link type",
		})
		return
	}

	c.JSON(http.StatusOK, linkType)
}

// Update handles PUT /api/v1/link-types/:id
func (h *LinkTypeHandler) Update(c *gin.Context) {
	// Parse ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid link type ID",
		})
		return
	}

	var input service.UpdateLinkTypeInput

	// Bind and validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate cardinality if provided
	if input.Cardinality != nil && !input.Cardinality.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cardinality value",
		})
		return
	}

	// Sanitize input to prevent XSS
	if input.DisplayName != nil {
		sanitized := validator.SanitizeString(*input.DisplayName)
		input.DisplayName = &sanitized
	}
	if input.Description != nil {
		sanitized := validator.SanitizeString(*input.Description)
		input.Description = &sanitized
	}
	if input.InverseDisplayName != nil {
		sanitized := validator.SanitizeString(*input.InverseDisplayName)
		input.InverseDisplayName = &sanitized
	}

	// Get user ID from context
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Update link type
	linkType, err := h.service.UpdateLinkType(c.Request.Context(), id, input, userID)
	if err != nil {
		if err == entity.ErrLinkTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Link type not found",
			})
			return
		}

		h.logger.Error("Failed to update link type",
			zap.String("id", id.String()),
			zap.String("user_id", userID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update link type",
		})
		return
	}

	c.JSON(http.StatusOK, linkType)
}

// Delete handles DELETE /api/v1/link-types/:id
func (h *LinkTypeHandler) Delete(c *gin.Context) {
	// Parse ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid link type ID",
		})
		return
	}

	// Get user ID from context
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Check if user has permission to delete
	if !middleware.HasRole(c, "admin") {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})
		return
	}

	// Delete link type
	err = h.service.DeleteLinkType(c.Request.Context(), id, userID)
	if err != nil {
		if err == entity.ErrLinkTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Link type not found",
			})
			return
		}

		h.logger.Error("Failed to delete link type",
			zap.String("id", id.String()),
			zap.String("user_id", userID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete link type",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetByObjectTypes handles GET /api/v1/link-types/by-object-types
func (h *LinkTypeHandler) GetByObjectTypes(c *gin.Context) {
	// Parse source and target IDs
	sourceIDStr := c.Query("source_id")
	targetIDStr := c.Query("target_id")

	if sourceIDStr == "" || targetIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both source_id and target_id are required",
		})
		return
	}

	sourceID, err := uuid.Parse(sourceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid source_id",
		})
		return
	}

	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid target_id",
		})
		return
	}

	// Get link types
	linkTypes, err := h.service.GetByObjectTypes(c.Request.Context(), sourceID, targetID)
	if err != nil {
		h.logger.Error("Failed to get link types by object types",
			zap.String("source_id", sourceID.String()),
			zap.String("target_id", targetID.String()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve link types",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": linkTypes,
	})
}

// ValidateCircularReference handles POST /api/v1/link-types/validate-circular
func (h *LinkTypeHandler) ValidateCircularReference(c *gin.Context) {
	var input struct {
		SourceObjectTypeID uuid.UUID `json:"sourceObjectTypeId" binding:"required"`
		TargetObjectTypeID uuid.UUID `json:"targetObjectTypeId" binding:"required"`
	}

	// Bind and validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Check circular reference
	hasCircular, err := h.service.CheckCircularReference(c.Request.Context(), input.SourceObjectTypeID, input.TargetObjectTypeID)
	if err != nil {
		h.logger.Error("Failed to check circular reference",
			zap.String("source_id", input.SourceObjectTypeID.String()),
			zap.String("target_id", input.TargetObjectTypeID.String()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to validate circular reference",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_circular_reference": hasCircular,
		"is_valid":               !hasCircular,
	})
}