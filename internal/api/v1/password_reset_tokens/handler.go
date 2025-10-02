package password_reset_tokens

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

// Handler handles password_reset_tokens requests
type Handler struct {
	passwordResetTokensRepo generated.PasswordresettokensRepository
	logger             *zap.Logger
}

// NewHandler creates a new password_reset_tokens handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		passwordResetTokensRepo: generated.NewPasswordresettokensRepository(db),
		logger:             logger,
	}
}

// CreatePasswordresettokens creates a new password_reset_tokens
// @Summary Create password_reset_tokens
// @Description Create a new password_reset_tokens record
// @Tags password_reset_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PasswordresettokensCreateRequest true "Create password_reset_tokens request"
// @Success 201 {object} PasswordresettokensResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /password_reset_tokens [post]
func (h *Handler) CreatePasswordresettokens(c *gin.Context) {
	var req PasswordresettokensCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	passwordResetTokens := &generated.Passwordresettokens{

		Userid: req.Userid,

		Token: req.Token,

		Expiresat: req.Expiresat,

		Used: req.Used,

		Usedat: req.Usedat,

	}

	ctx := context.Background()
	if err := h.passwordResetTokensRepo.Create(ctx, passwordResetTokens); err != nil {
		h.logger.Error("Failed to create password_reset_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create password_reset_tokens")
		return
	}

	response := PasswordresettokensResponse{

		Id: passwordResetTokens.Id,

		Userid: passwordResetTokens.Userid,

		Token: passwordResetTokens.Token,

		Expiresat: passwordResetTokens.Expiresat,

		Used: passwordResetTokens.Used,

		Usedat: passwordResetTokens.Usedat,

		Createdat: passwordResetTokens.Createdat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "password_reset_tokens created successfully", response)
}

// GetPasswordresettokens gets a password_reset_tokens by ID
// @Summary Get password_reset_tokens
// @Description Get a password_reset_tokens by ID
// @Tags password_reset_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "password_reset_tokens ID"
// @Success 200 {object} PasswordresettokensResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /password_reset_tokens/{id} [get]
func (h *Handler) GetPasswordresettokens(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	passwordResetTokens, err := h.passwordResetTokensRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "password_reset_tokens not found")
			return
		}
		h.logger.Error("Failed to get password_reset_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get password_reset_tokens")
		return
	}

	response := PasswordresettokensResponse{

		Id: passwordResetTokens.Id,

		Userid: passwordResetTokens.Userid,

		Token: passwordResetTokens.Token,

		Expiresat: passwordResetTokens.Expiresat,

		Used: passwordResetTokens.Used,

		Usedat: passwordResetTokens.Usedat,

		Createdat: passwordResetTokens.Createdat,

	}

	utils.SuccessResponse(c, http.StatusOK, "password_reset_tokens retrieved successfully", response)
}

// GetAllPasswordresettokenss gets all password_reset_tokenss with pagination
// @Summary Get all password_reset_tokenss
// @Description Get all password_reset_tokenss with pagination
// @Tags password_reset_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /password_reset_tokens [get]
func (h *Handler) GetAllPasswordresettokenss(c *gin.Context) {
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
	passwordResetTokenss, total, err := h.passwordResetTokensRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get password_reset_tokenss", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get password_reset_tokenss")
		return
	}

	var responses []PasswordresettokensResponse
	for _, passwordResetTokens := range passwordResetTokenss {
		responses = append(responses, PasswordresettokensResponse{

			Id: passwordResetTokens.Id,

			Userid: passwordResetTokens.Userid,

			Token: passwordResetTokens.Token,

			Expiresat: passwordResetTokens.Expiresat,

			Used: passwordResetTokens.Used,

			Usedat: passwordResetTokens.Usedat,

			Createdat: passwordResetTokens.Createdat,

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

	utils.SuccessResponse(c, http.StatusOK, "password_reset_tokenss retrieved successfully", pagination)
}

// UpdatePasswordresettokens updates a password_reset_tokens
// @Summary Update password_reset_tokens
// @Description Update a password_reset_tokens by ID
// @Tags password_reset_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "password_reset_tokens ID"
// @Param request body PasswordresettokensUpdateRequest true "Update password_reset_tokens request"
// @Success 200 {object} PasswordresettokensResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /password_reset_tokens/{id} [put]
func (h *Handler) UpdatePasswordresettokens(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req PasswordresettokensUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	passwordResetTokens, err := h.passwordResetTokensRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "password_reset_tokens not found")
			return
		}
		h.logger.Error("Failed to get password_reset_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get password_reset_tokens")
		return
	}


	
	passwordResetTokens.Userid = req.Userid
	

	
	passwordResetTokens.Token = req.Token
	

	
	passwordResetTokens.Expiresat = req.Expiresat
	

	
	passwordResetTokens.Used = req.Used
	

	
	passwordResetTokens.Usedat = req.Usedat
	


	if err := h.passwordResetTokensRepo.Update(ctx, passwordResetTokens); err != nil {
		h.logger.Error("Failed to update password_reset_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update password_reset_tokens")
		return
	}

	response := PasswordresettokensResponse{

		Id: passwordResetTokens.Id,

		Userid: passwordResetTokens.Userid,

		Token: passwordResetTokens.Token,

		Expiresat: passwordResetTokens.Expiresat,

		Used: passwordResetTokens.Used,

		Usedat: passwordResetTokens.Usedat,

		Createdat: passwordResetTokens.Createdat,

	}

	utils.SuccessResponse(c, http.StatusOK, "password_reset_tokens updated successfully", response)
}

// DeletePasswordresettokens deletes a password_reset_tokens
// @Summary Delete password_reset_tokens
// @Description Delete a password_reset_tokens by ID
// @Tags password_reset_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "password_reset_tokens ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /password_reset_tokens/{id} [delete]
func (h *Handler) DeletePasswordresettokens(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.passwordResetTokensRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "password_reset_tokens not found")
			return
		}
		h.logger.Error("Failed to delete password_reset_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete password_reset_tokens")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "password_reset_tokens deleted successfully", nil)
}
