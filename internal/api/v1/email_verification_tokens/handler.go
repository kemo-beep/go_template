package email_verification_tokens

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

// Handler handles email_verification_tokens requests
type Handler struct {
	emailVerificationTokensRepo generated.EmailverificationtokensRepository
	logger             *zap.Logger
}

// NewHandler creates a new email_verification_tokens handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		emailVerificationTokensRepo: generated.NewEmailverificationtokensRepository(db),
		logger:             logger,
	}
}

// CreateEmailverificationtokens creates a new email_verification_tokens
// @Summary Create email_verification_tokens
// @Description Create a new email_verification_tokens record
// @Tags email_verification_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body EmailverificationtokensCreateRequest true "Create email_verification_tokens request"
// @Success 201 {object} EmailverificationtokensResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /email_verification_tokens [post]
func (h *Handler) CreateEmailverificationtokens(c *gin.Context) {
	var req EmailverificationtokensCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	emailVerificationTokens := &generated.Emailverificationtokens{

		Userid: req.Userid,

		Email: req.Email,

		Token: req.Token,

		Expiresat: req.Expiresat,

		Used: req.Used,

		Usedat: req.Usedat,

	}

	ctx := context.Background()
	if err := h.emailVerificationTokensRepo.Create(ctx, emailVerificationTokens); err != nil {
		h.logger.Error("Failed to create email_verification_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create email_verification_tokens")
		return
	}

	response := EmailverificationtokensResponse{

		Id: emailVerificationTokens.Id,

		Userid: emailVerificationTokens.Userid,

		Email: emailVerificationTokens.Email,

		Token: emailVerificationTokens.Token,

		Expiresat: emailVerificationTokens.Expiresat,

		Used: emailVerificationTokens.Used,

		Usedat: emailVerificationTokens.Usedat,

		Createdat: emailVerificationTokens.Createdat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "email_verification_tokens created successfully", response)
}

// GetEmailverificationtokens gets a email_verification_tokens by ID
// @Summary Get email_verification_tokens
// @Description Get a email_verification_tokens by ID
// @Tags email_verification_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "email_verification_tokens ID"
// @Success 200 {object} EmailverificationtokensResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /email_verification_tokens/{id} [get]
func (h *Handler) GetEmailverificationtokens(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	emailVerificationTokens, err := h.emailVerificationTokensRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "email_verification_tokens not found")
			return
		}
		h.logger.Error("Failed to get email_verification_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get email_verification_tokens")
		return
	}

	response := EmailverificationtokensResponse{

		Id: emailVerificationTokens.Id,

		Userid: emailVerificationTokens.Userid,

		Email: emailVerificationTokens.Email,

		Token: emailVerificationTokens.Token,

		Expiresat: emailVerificationTokens.Expiresat,

		Used: emailVerificationTokens.Used,

		Usedat: emailVerificationTokens.Usedat,

		Createdat: emailVerificationTokens.Createdat,

	}

	utils.SuccessResponse(c, http.StatusOK, "email_verification_tokens retrieved successfully", response)
}

// GetAllEmailverificationtokenss gets all email_verification_tokenss with pagination
// @Summary Get all email_verification_tokenss
// @Description Get all email_verification_tokenss with pagination
// @Tags email_verification_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /email_verification_tokens [get]
func (h *Handler) GetAllEmailverificationtokenss(c *gin.Context) {
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
	emailVerificationTokenss, total, err := h.emailVerificationTokensRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get email_verification_tokenss", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get email_verification_tokenss")
		return
	}

	var responses []EmailverificationtokensResponse
	for _, emailVerificationTokens := range emailVerificationTokenss {
		responses = append(responses, EmailverificationtokensResponse{

			Id: emailVerificationTokens.Id,

			Userid: emailVerificationTokens.Userid,

			Email: emailVerificationTokens.Email,

			Token: emailVerificationTokens.Token,

			Expiresat: emailVerificationTokens.Expiresat,

			Used: emailVerificationTokens.Used,

			Usedat: emailVerificationTokens.Usedat,

			Createdat: emailVerificationTokens.Createdat,

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

	utils.SuccessResponse(c, http.StatusOK, "email_verification_tokenss retrieved successfully", pagination)
}

// UpdateEmailverificationtokens updates a email_verification_tokens
// @Summary Update email_verification_tokens
// @Description Update a email_verification_tokens by ID
// @Tags email_verification_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "email_verification_tokens ID"
// @Param request body EmailverificationtokensUpdateRequest true "Update email_verification_tokens request"
// @Success 200 {object} EmailverificationtokensResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /email_verification_tokens/{id} [put]
func (h *Handler) UpdateEmailverificationtokens(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req EmailverificationtokensUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	emailVerificationTokens, err := h.emailVerificationTokensRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "email_verification_tokens not found")
			return
		}
		h.logger.Error("Failed to get email_verification_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get email_verification_tokens")
		return
	}


	
	emailVerificationTokens.Userid = req.Userid
	

	
	emailVerificationTokens.Email = req.Email
	

	
	emailVerificationTokens.Token = req.Token
	

	
	emailVerificationTokens.Expiresat = req.Expiresat
	

	
	emailVerificationTokens.Used = req.Used
	

	
	emailVerificationTokens.Usedat = req.Usedat
	


	if err := h.emailVerificationTokensRepo.Update(ctx, emailVerificationTokens); err != nil {
		h.logger.Error("Failed to update email_verification_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update email_verification_tokens")
		return
	}

	response := EmailverificationtokensResponse{

		Id: emailVerificationTokens.Id,

		Userid: emailVerificationTokens.Userid,

		Email: emailVerificationTokens.Email,

		Token: emailVerificationTokens.Token,

		Expiresat: emailVerificationTokens.Expiresat,

		Used: emailVerificationTokens.Used,

		Usedat: emailVerificationTokens.Usedat,

		Createdat: emailVerificationTokens.Createdat,

	}

	utils.SuccessResponse(c, http.StatusOK, "email_verification_tokens updated successfully", response)
}

// DeleteEmailverificationtokens deletes a email_verification_tokens
// @Summary Delete email_verification_tokens
// @Description Delete a email_verification_tokens by ID
// @Tags email_verification_tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "email_verification_tokens ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /email_verification_tokens/{id} [delete]
func (h *Handler) DeleteEmailverificationtokens(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.emailVerificationTokensRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "email_verification_tokens not found")
			return
		}
		h.logger.Error("Failed to delete email_verification_tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete email_verification_tokens")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "email_verification_tokens deleted successfully", nil)
}
