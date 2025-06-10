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

// ObjectTypeHandler handles object type related requests
type ObjectTypeHandler struct {
	service *service.ObjectTypeService
	logger  *zap.Logger
}

// NewObjectTypeHandler creates a new object type handler
func NewObjectTypeHandler(service *service.ObjectTypeService, logger *zap.Logger) *ObjectTypeHandler {
	return &ObjectTypeHandler{
		service: service,
		logger:  logger,
	}
}

// List handles GET /api/v1/object-types
func (h *ObjectTypeHandler) List(c *gin.Context) {
	// Parse query parameters
	filter := repository.ObjectTypeFilter{
		PageSize: 20, // Default page size
	}

	// Parse category filter
	if category := c.Query("category"); category != "" {
		filter.Category = &category
	}

	// Parse tags filter
	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filter.Tags = tags
	}

	// Parse pagination
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			filter.PageSize = pageSize
		}
	}

	if cursor := c.Query("cursor"); cursor != "" {
		filter.PageCursor = cursor
	}

	// Parse sort
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filter.SortOrder = sortOrder
	}

	// Get object types
	objectTypes, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to list object types", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve object types",
		})
		return
	}

	// Generate next cursor if needed
	var nextCursor string
	if len(objectTypes) == filter.PageSize {
		lastItem := objectTypes[len(objectTypes)-1]
		nextCursor = encodeCursor(lastItem.CreatedAt, lastItem.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": objectTypes,
		"pagination": gin.H{
			"next_cursor": nextCursor,
			"page_size":   filter.PageSize,
		},
	})
}

// Create handles POST /api/v1/object-types
func (h *ObjectTypeHandler) Create(c *gin.Context) {
	var input service.CreateObjectTypeInput

	// Bind and validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Additional validation
	if err := validator.ValidateObjectTypeName(input.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid object type name",
			"details": err.Error(),
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

	// Get user ID from context
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Create object type
	objectType, err := h.service.CreateObjectType(c.Request.Context(), input, userID)
	if err != nil {
		h.logger.Error("Failed to create object type", 
			zap.String("user_id", userID),
			zap.String("name", input.Name),
			zap.Error(err))

		// Handle specific errors
		switch err {
		case entity.ErrObjectTypeNameExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": "Object type name already exists",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create object type",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, objectType)
}

// Get handles GET /api/v1/object-types/:id
func (h *ObjectTypeHandler) Get(c *gin.Context) {
	// Parse ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid object type ID",
		})
		return
	}

	// Get object type
	objectType, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == entity.ErrObjectTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Object type not found",
			})
			return
		}

		h.logger.Error("Failed to get object type", 
			zap.String("id", id.String()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve object type",
		})
		return
	}

	c.JSON(http.StatusOK, objectType)
}

// Update handles PUT /api/v1/object-types/:id
func (h *ObjectTypeHandler) Update(c *gin.Context) {
	// Parse ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid object type ID",
		})
		return
	}

	var input service.UpdateObjectTypeInput

	// Bind and validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
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

	// Get user ID from context
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Update object type
	objectType, err := h.service.UpdateObjectType(c.Request.Context(), id, input, userID)
	if err != nil {
		if err == entity.ErrObjectTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Object type not found",
			})
			return
		}

		h.logger.Error("Failed to update object type", 
			zap.String("id", id.String()),
			zap.String("user_id", userID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update object type",
		})
		return
	}

	c.JSON(http.StatusOK, objectType)
}

// Delete handles DELETE /api/v1/object-types/:id
func (h *ObjectTypeHandler) Delete(c *gin.Context) {
	// Parse ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid object type ID",
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

	// Delete object type
	err = h.service.DeleteObjectType(c.Request.Context(), id, userID)
	if err != nil {
		if err == entity.ErrObjectTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Object type not found",
			})
			return
		}

		h.logger.Error("Failed to delete object type", 
			zap.String("id", id.String()),
			zap.String("user_id", userID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete object type",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Search handles GET /api/v1/search
func (h *ObjectTypeHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
		})
		return
	}

	// Sanitize query
	query = validator.SanitizeString(query)

	// Parse limit
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	// Search object types
	results, err := h.service.Search(c.Request.Context(), query, limit)
	if err != nil {
		h.logger.Error("Failed to search object types", 
			zap.String("query", query),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query": query,
		"results": results,
		"count": len(results),
	})
}

// CompareVersions handles GET /api/v1/object-types/:id/versions/compare
func (h *ObjectTypeHandler) CompareVersions(c *gin.Context) {
	// Parse ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid object type ID",
		})
		return
	}

	// Parse version parameters
	v1Str := c.Query("v1")
	v2Str := c.Query("v2")

	if v1Str == "" || v2Str == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both v1 and v2 version parameters are required",
		})
		return
	}

	v1, err := strconv.Atoi(v1Str)
	if err != nil || v1 < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid v1 version number",
		})
		return
	}

	v2, err := strconv.Atoi(v2Str)
	if err != nil || v2 < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid v2 version number",
		})
		return
	}

	// Compare versions
	diff, err := h.service.CompareVersions(c.Request.Context(), id, v1, v2)
	if err != nil {
		if err == entity.ErrObjectTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Object type not found",
			})
			return
		}

		h.logger.Error("Failed to compare versions", 
			zap.String("id", id.String()),
			zap.Int("v1", v1),
			zap.Int("v2", v2),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to compare versions",
		})
		return
	}

	c.JSON(http.StatusOK, diff)
}

// Helper function to encode cursor
func encodeCursor(timestamp time.Time, id uuid.UUID) string {
	// This should match the implementation in the repository
	data := fmt.Sprintf("%d:%s", timestamp.Unix(), id.String())
	return base64.StdEncoding.EncodeToString([]byte(data))
}