package refresh_tokens

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

// Handler handles refresh_tokens requests
type Handler struct {
	refreshTokensRepo generated.RefreshtokensRepository
	logger             *zap.Logger
}

// NewHandler creates a new refresh_tokens handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		refreshTokensRepo: generated.NewRefreshtokensRepository(db),
		logger:             logger,
	}
}

// CreateRefreshtokens creates a new refresh_tokens
// @Summary Create refresh_tokens
// @Description Create a new refresh_tokens record
// @Tags refresh_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RefreshtokensCreateRequest true "Create refresh_tokens request"
// @Success 201 {object} RefreshtokensResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /refresh_tokens [post]
func (h *Handler) CreateRefreshtokens(c *gin.Context) {
	var req RefreshtokensCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	refreshTokens := &generated.Refreshtokens{

		Userid: req.Userid,

		Token: req.Token,

		Expiresat: req.Expiresat,

		Isrevoked: req.Isrevoked,

		Deletedat: req.Deletedat,

	}

	ctx := context.Background()
	if err := h.refreshTokensRepo.Create(ctx, refreshTokens); err != nil {
		h.logger.Error("Failed to create refresh_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create refresh_tokens")
		return
	}

	response := RefreshtokensResponse{

		Id: refreshTokens.Id,

		Userid: refreshTokens.Userid,

		Token: refreshTokens.Token,

		Expiresat: refreshTokens.Expiresat,

		Isrevoked: refreshTokens.Isrevoked,

		Createdat: refreshTokens.Createdat,

		Updatedat: refreshTokens.Updatedat,

		Deletedat: refreshTokens.Deletedat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "refresh_tokens created successfully", response)
}

// GetRefreshtokens gets a refresh_tokens by ID
// @Summary Get refresh_tokens
// @Description Get a refresh_tokens by ID
// @Tags refresh_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "refresh_tokens ID"
// @Success 200 {object} RefreshtokensResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /refresh_tokens/{id} [get]
func (h *Handler) GetRefreshtokens(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	refreshTokens, err := h.refreshTokensRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "refresh_tokens not found")
			return
		}
		h.logger.Error("Failed to get refresh_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get refresh_tokens")
		return
	}

	response := RefreshtokensResponse{

		Id: refreshTokens.Id,

		Userid: refreshTokens.Userid,

		Token: refreshTokens.Token,

		Expiresat: refreshTokens.Expiresat,

		Isrevoked: refreshTokens.Isrevoked,

		Createdat: refreshTokens.Createdat,

		Updatedat: refreshTokens.Updatedat,

		Deletedat: refreshTokens.Deletedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "refresh_tokens retrieved successfully", response)
}

// GetAllRefreshtokenss gets all refresh_tokenss with pagination
// @Summary Get all refresh_tokenss
// @Description Get all refresh_tokenss with pagination
// @Tags refresh_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /refresh_tokens [get]
func (h *Handler) GetAllRefreshtokenss(c *gin.Context) {
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
	refreshTokenss, total, err := h.refreshTokensRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get refresh_tokenss", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get refresh_tokenss")
		return
	}

	var responses []RefreshtokensResponse
	for _, refreshTokens := range refreshTokenss {
		responses = append(responses, RefreshtokensResponse{

			Id: refreshTokens.Id,

			Userid: refreshTokens.Userid,

			Token: refreshTokens.Token,

			Expiresat: refreshTokens.Expiresat,

			Isrevoked: refreshTokens.Isrevoked,

			Createdat: refreshTokens.Createdat,

			Updatedat: refreshTokens.Updatedat,

			Deletedat: refreshTokens.Deletedat,

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

	utils.SuccessResponse(c, http.StatusOK, "refresh_tokenss retrieved successfully", pagination)
}

// UpdateRefreshtokens updates a refresh_tokens
// @Summary Update refresh_tokens
// @Description Update a refresh_tokens by ID
// @Tags refresh_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "refresh_tokens ID"
// @Param request body RefreshtokensUpdateRequest true "Update refresh_tokens request"
// @Success 200 {object} RefreshtokensResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /refresh_tokens/{id} [put]
func (h *Handler) UpdateRefreshtokens(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req RefreshtokensUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	refreshTokens, err := h.refreshTokensRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "refresh_tokens not found")
			return
		}
		h.logger.Error("Failed to get refresh_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get refresh_tokens")
		return
	}


	
	refreshTokens.Userid = req.Userid
	

	
	refreshTokens.Token = req.Token
	

	
	refreshTokens.Expiresat = req.Expiresat
	

	
	refreshTokens.Isrevoked = req.Isrevoked
	

	
	refreshTokens.Deletedat = req.Deletedat
	


	if err := h.refreshTokensRepo.Update(ctx, refreshTokens); err != nil {
		h.logger.Error("Failed to update refresh_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update refresh_tokens")
		return
	}

	response := RefreshtokensResponse{

		Id: refreshTokens.Id,

		Userid: refreshTokens.Userid,

		Token: refreshTokens.Token,

		Expiresat: refreshTokens.Expiresat,

		Isrevoked: refreshTokens.Isrevoked,

		Createdat: refreshTokens.Createdat,

		Updatedat: refreshTokens.Updatedat,

		Deletedat: refreshTokens.Deletedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "refresh_tokens updated successfully", response)
}

// DeleteRefreshtokens deletes a refresh_tokens
// @Summary Delete refresh_tokens
// @Description Delete a refresh_tokens by ID
// @Tags refresh_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "refresh_tokens ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /refresh_tokens/{id} [delete]
func (h *Handler) DeleteRefreshtokens(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.refreshTokensRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "refresh_tokens not found")
			return
		}
		h.logger.Error("Failed to delete refresh_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete refresh_tokens")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "refresh_tokens deleted successfully", nil)
}
