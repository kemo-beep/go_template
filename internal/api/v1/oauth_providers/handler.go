package oauth_providers

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

// Handler handles oauth_providers requests
type Handler struct {
	oauthProvidersRepo generated.OauthprovidersRepository
	logger             *zap.Logger
}

// NewHandler creates a new oauth_providers handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		oauthProvidersRepo: generated.NewOauthprovidersRepository(db),
		logger:             logger,
	}
}

// CreateOauthproviders creates a new oauth_providers
// @Summary Create oauth_providers
// @Description Create a new oauth_providers record
// @Tags oauth_providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body OauthprovidersCreateRequest true "Create oauth_providers request"
// @Success 201 {object} OauthprovidersResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /oauth_providers [post]
func (h *Handler) CreateOauthproviders(c *gin.Context) {
	var req OauthprovidersCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	oauthProviders := &generated.Oauthproviders{

		Userid: req.Userid,

		Provider: req.Provider,

		Provideruserid: req.Provideruserid,

		Accesstoken: req.Accesstoken,

		Refreshtoken: req.Refreshtoken,

		Tokenexpiresat: req.Tokenexpiresat,

		Profiledata: req.Profiledata,

	}

	ctx := context.Background()
	if err := h.oauthProvidersRepo.Create(ctx, oauthProviders); err != nil {
		h.logger.Error("Failed to create oauth_providers", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create oauth_providers")
		return
	}

	response := OauthprovidersResponse{

		Id: oauthProviders.Id,

		Userid: oauthProviders.Userid,

		Provider: oauthProviders.Provider,

		Provideruserid: oauthProviders.Provideruserid,

		Accesstoken: oauthProviders.Accesstoken,

		Refreshtoken: oauthProviders.Refreshtoken,

		Tokenexpiresat: oauthProviders.Tokenexpiresat,

		Profiledata: oauthProviders.Profiledata,

		Createdat: oauthProviders.Createdat,

		Updatedat: oauthProviders.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "oauth_providers created successfully", response)
}

// GetOauthproviders gets a oauth_providers by ID
// @Summary Get oauth_providers
// @Description Get a oauth_providers by ID
// @Tags oauth_providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "oauth_providers ID"
// @Success 200 {object} OauthprovidersResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /oauth_providers/{id} [get]
func (h *Handler) GetOauthproviders(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	oauthProviders, err := h.oauthProvidersRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "oauth_providers not found")
			return
		}
		h.logger.Error("Failed to get oauth_providers", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get oauth_providers")
		return
	}

	response := OauthprovidersResponse{

		Id: oauthProviders.Id,

		Userid: oauthProviders.Userid,

		Provider: oauthProviders.Provider,

		Provideruserid: oauthProviders.Provideruserid,

		Accesstoken: oauthProviders.Accesstoken,

		Refreshtoken: oauthProviders.Refreshtoken,

		Tokenexpiresat: oauthProviders.Tokenexpiresat,

		Profiledata: oauthProviders.Profiledata,

		Createdat: oauthProviders.Createdat,

		Updatedat: oauthProviders.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "oauth_providers retrieved successfully", response)
}

// GetAllOauthproviderss gets all oauth_providerss with pagination
// @Summary Get all oauth_providerss
// @Description Get all oauth_providerss with pagination
// @Tags oauth_providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /oauth_providers [get]
func (h *Handler) GetAllOauthproviderss(c *gin.Context) {
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
	oauthProviderss, total, err := h.oauthProvidersRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get oauth_providerss", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get oauth_providerss")
		return
	}

	var responses []OauthprovidersResponse
	for _, oauthProviders := range oauthProviderss {
		responses = append(responses, OauthprovidersResponse{

			Id: oauthProviders.Id,

			Userid: oauthProviders.Userid,

			Provider: oauthProviders.Provider,

			Provideruserid: oauthProviders.Provideruserid,

			Accesstoken: oauthProviders.Accesstoken,

			Refreshtoken: oauthProviders.Refreshtoken,

			Tokenexpiresat: oauthProviders.Tokenexpiresat,

			Profiledata: oauthProviders.Profiledata,

			Createdat: oauthProviders.Createdat,

			Updatedat: oauthProviders.Updatedat,

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

	utils.SuccessResponse(c, http.StatusOK, "oauth_providerss retrieved successfully", pagination)
}

// UpdateOauthproviders updates a oauth_providers
// @Summary Update oauth_providers
// @Description Update a oauth_providers by ID
// @Tags oauth_providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "oauth_providers ID"
// @Param request body OauthprovidersUpdateRequest true "Update oauth_providers request"
// @Success 200 {object} OauthprovidersResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /oauth_providers/{id} [put]
func (h *Handler) UpdateOauthproviders(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req OauthprovidersUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	oauthProviders, err := h.oauthProvidersRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "oauth_providers not found")
			return
		}
		h.logger.Error("Failed to get oauth_providers", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get oauth_providers")
		return
	}


	
	oauthProviders.Userid = req.Userid
	

	
	oauthProviders.Provider = req.Provider
	

	
	oauthProviders.Provideruserid = req.Provideruserid
	

	
	oauthProviders.Accesstoken = req.Accesstoken
	

	
	oauthProviders.Refreshtoken = req.Refreshtoken
	

	
	oauthProviders.Tokenexpiresat = req.Tokenexpiresat
	

	
	oauthProviders.Profiledata = req.Profiledata
	


	if err := h.oauthProvidersRepo.Update(ctx, oauthProviders); err != nil {
		h.logger.Error("Failed to update oauth_providers", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update oauth_providers")
		return
	}

	response := OauthprovidersResponse{

		Id: oauthProviders.Id,

		Userid: oauthProviders.Userid,

		Provider: oauthProviders.Provider,

		Provideruserid: oauthProviders.Provideruserid,

		Accesstoken: oauthProviders.Accesstoken,

		Refreshtoken: oauthProviders.Refreshtoken,

		Tokenexpiresat: oauthProviders.Tokenexpiresat,

		Profiledata: oauthProviders.Profiledata,

		Createdat: oauthProviders.Createdat,

		Updatedat: oauthProviders.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "oauth_providers updated successfully", response)
}

// DeleteOauthproviders deletes a oauth_providers
// @Summary Delete oauth_providers
// @Description Delete a oauth_providers by ID
// @Tags oauth_providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "oauth_providers ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /oauth_providers/{id} [delete]
func (h *Handler) DeleteOauthproviders(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.oauthProvidersRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "oauth_providers not found")
			return
		}
		h.logger.Error("Failed to delete oauth_providers", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete oauth_providers")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "oauth_providers deleted successfully", nil)
}
