package api_keys

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/db/repository/generated"
	"go-mobile-backend-template/internal/utils"
)

// Handler handles api_keys requests
type Handler struct {
	apiKeysRepo generated.ApikeysRepository
	logger             *zap.Logger
}

// NewHandler creates a new api_keys handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		apiKeysRepo: generated.NewApikeysRepository(db),
		logger:             logger,
	}
}

// CreateApikeys creates a new api_keys
// @Summary Create api_keys
// @Description Create a new api_keys record
// @Tags api_keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ApikeysCreateRequest true "Create api_keys request"
// @Success 201 {object} ApikeysResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api_keys [post]
func (h *Handler) CreateApikeys(c *gin.Context) {
	var req ApikeysCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	apiKeys := &generated.Apikeys{

		Userid: req.Userid,

		Name: req.Name,

		Keyhash: req.Keyhash,

		Prefix: req.Prefix,

		Scopes: req.Scopes,

		Ratelimit: req.Ratelimit,

		Isactive: req.Isactive,

		Lastusedat: req.Lastusedat,

		Expiresat: req.Expiresat,

	}

	ctx := context.Background()
	if err := h.apiKeysRepo.Create(ctx, apiKeys); err != nil {
		h.logger.Error("Failed to create api_keys", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create api_keys")
		return
	}

	response := ApikeysResponse{

		Id: apiKeys.Id,

		Userid: apiKeys.Userid,

		Name: apiKeys.Name,

		Keyhash: apiKeys.Keyhash,

		Prefix: apiKeys.Prefix,

		Scopes: apiKeys.Scopes,

		Ratelimit: apiKeys.Ratelimit,

		Isactive: apiKeys.Isactive,

		Lastusedat: apiKeys.Lastusedat,

		Expiresat: apiKeys.Expiresat,

		Createdat: apiKeys.Createdat,

		Updatedat: apiKeys.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "api_keys created successfully", response)
}

// GetApikeys gets a api_keys by ID
// @Summary Get api_keys
// @Description Get a api_keys by ID
// @Tags api_keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "api_keys ID"
// @Success 200 {object} ApikeysResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api_keys/{id} [get]
func (h *Handler) GetApikeys(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	apiKeys, err := h.apiKeysRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "api_keys not found")
			return
		}
		h.logger.Error("Failed to get api_keys", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get api_keys")
		return
	}

	response := ApikeysResponse{

		Id: apiKeys.Id,

		Userid: apiKeys.Userid,

		Name: apiKeys.Name,

		Keyhash: apiKeys.Keyhash,

		Prefix: apiKeys.Prefix,

		Scopes: apiKeys.Scopes,

		Ratelimit: apiKeys.Ratelimit,

		Isactive: apiKeys.Isactive,

		Lastusedat: apiKeys.Lastusedat,

		Expiresat: apiKeys.Expiresat,

		Createdat: apiKeys.Createdat,

		Updatedat: apiKeys.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "api_keys retrieved successfully", response)
}

// GetAllApikeyss gets all api_keyss with pagination
// @Summary Get all api_keyss
// @Description Get all api_keyss with pagination
// @Tags api_keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api_keys [get]
func (h *Handler) GetAllApikeyss(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	ctx := context.Background()
	apiKeyss, total, err := h.apiKeysRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get api_keyss", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get api_keyss")
		return
	}

	var responses []ApikeysResponse
	for _, apiKeys := range apiKeyss {
		responses = append(responses, ApikeysResponse{

			Id: apiKeys.Id,

			Userid: apiKeys.Userid,

			Name: apiKeys.Name,

			Keyhash: apiKeys.Keyhash,

			Prefix: apiKeys.Prefix,

			Scopes: apiKeys.Scopes,

			Ratelimit: apiKeys.Ratelimit,

			Isactive: apiKeys.Isactive,

			Lastusedat: apiKeys.Lastusedat,

			Expiresat: apiKeys.Expiresat,

			Createdat: apiKeys.Createdat,

			Updatedat: apiKeys.Updatedat,

		})
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	pagination := PaginationResponse{
		Data: responses,
		Pagination: PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "api_keyss retrieved successfully", pagination)
}

// UpdateApikeys updates a api_keys
// @Summary Update api_keys
// @Description Update a api_keys by ID
// @Tags api_keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "api_keys ID"
// @Param request body ApikeysUpdateRequest true "Update api_keys request"
// @Success 200 {object} ApikeysResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api_keys/{id} [put]
func (h *Handler) UpdateApikeys(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req ApikeysUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	apiKeys, err := h.apiKeysRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "api_keys not found")
			return
		}
		h.logger.Error("Failed to get api_keys", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get api_keys")
		return
	}


	
	apiKeys.Userid = req.Userid
	

	
	apiKeys.Name = req.Name
	

	
	apiKeys.Keyhash = req.Keyhash
	

	
	apiKeys.Prefix = req.Prefix
	

	
	apiKeys.Scopes = req.Scopes
	

	
	apiKeys.Ratelimit = req.Ratelimit
	

	
	apiKeys.Isactive = req.Isactive
	

	
	apiKeys.Lastusedat = req.Lastusedat
	

	
	apiKeys.Expiresat = req.Expiresat
	


	if err := h.apiKeysRepo.Update(ctx, apiKeys); err != nil {
		h.logger.Error("Failed to update api_keys", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update api_keys")
		return
	}

	response := ApikeysResponse{

		Id: apiKeys.Id,

		Userid: apiKeys.Userid,

		Name: apiKeys.Name,

		Keyhash: apiKeys.Keyhash,

		Prefix: apiKeys.Prefix,

		Scopes: apiKeys.Scopes,

		Ratelimit: apiKeys.Ratelimit,

		Isactive: apiKeys.Isactive,

		Lastusedat: apiKeys.Lastusedat,

		Expiresat: apiKeys.Expiresat,

		Createdat: apiKeys.Createdat,

		Updatedat: apiKeys.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "api_keys updated successfully", response)
}

// DeleteApikeys deletes a api_keys
// @Summary Delete api_keys
// @Description Delete a api_keys by ID
// @Tags api_keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "api_keys ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api_keys/{id} [delete]
func (h *Handler) DeleteApikeys(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.apiKeysRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "api_keys not found")
			return
		}
		h.logger.Error("Failed to delete api_keys", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete api_keys")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "api_keys deleted successfully", nil)
}
