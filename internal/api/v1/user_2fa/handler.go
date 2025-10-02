package user_2fa

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

// Handler handles user_2fa requests
type Handler struct {
	user2faRepo generated.User2faRepository
	logger             *zap.Logger
}

// NewHandler creates a new user_2fa handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		user2faRepo: generated.NewUser2faRepository(db),
		logger:             logger,
	}
}

// CreateUser2fa creates a new user_2fa
// @Summary Create user_2fa
// @Description Create a new user_2fa record
// @Tags user_2fa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body User2faCreateRequest true "Create user_2fa request"
// @Success 201 {object} User2faResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /user_2fa [post]
func (h *Handler) CreateUser2fa(c *gin.Context) {
	var req User2faCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	user2fa := &generated.User2fa{

		Userid: req.Userid,

		Secret: req.Secret,

		Backupcodes: req.Backupcodes,

		Isenabled: req.Isenabled,

		Enabledat: req.Enabledat,

		Lastusedat: req.Lastusedat,

	}

	ctx := context.Background()
	if err := h.user2faRepo.Create(ctx, user2fa); err != nil {
		h.logger.Error("Failed to create user_2fa", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user_2fa")
		return
	}

	response := User2faResponse{

		Id: user2fa.Id,

		Userid: user2fa.Userid,

		Secret: user2fa.Secret,

		Backupcodes: user2fa.Backupcodes,

		Isenabled: user2fa.Isenabled,

		Enabledat: user2fa.Enabledat,

		Lastusedat: user2fa.Lastusedat,

		Createdat: user2fa.Createdat,

		Updatedat: user2fa.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "user_2fa created successfully", response)
}

// GetUser2fa gets a user_2fa by ID
// @Summary Get user_2fa
// @Description Get a user_2fa by ID
// @Tags user_2fa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "user_2fa ID"
// @Success 200 {object} User2faResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /user_2fa/{id} [get]
func (h *Handler) GetUser2fa(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	user2fa, err := h.user2faRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "user_2fa not found")
			return
		}
		h.logger.Error("Failed to get user_2fa", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user_2fa")
		return
	}

	response := User2faResponse{

		Id: user2fa.Id,

		Userid: user2fa.Userid,

		Secret: user2fa.Secret,

		Backupcodes: user2fa.Backupcodes,

		Isenabled: user2fa.Isenabled,

		Enabledat: user2fa.Enabledat,

		Lastusedat: user2fa.Lastusedat,

		Createdat: user2fa.Createdat,

		Updatedat: user2fa.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "user_2fa retrieved successfully", response)
}

// GetAllUser2fas gets all user_2fas with pagination
// @Summary Get all user_2fas
// @Description Get all user_2fas with pagination
// @Tags user_2fa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /user_2fa [get]
func (h *Handler) GetAllUser2fas(c *gin.Context) {
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
	user2fas, total, err := h.user2faRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get user_2fas", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user_2fas")
		return
	}

	var responses []User2faResponse
	for _, user2fa := range user2fas {
		responses = append(responses, User2faResponse{

			Id: user2fa.Id,

			Userid: user2fa.Userid,

			Secret: user2fa.Secret,

			Backupcodes: user2fa.Backupcodes,

			Isenabled: user2fa.Isenabled,

			Enabledat: user2fa.Enabledat,

			Lastusedat: user2fa.Lastusedat,

			Createdat: user2fa.Createdat,

			Updatedat: user2fa.Updatedat,

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

	utils.SuccessResponse(c, http.StatusOK, "user_2fas retrieved successfully", pagination)
}

// UpdateUser2fa updates a user_2fa
// @Summary Update user_2fa
// @Description Update a user_2fa by ID
// @Tags user_2fa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "user_2fa ID"
// @Param request body User2faUpdateRequest true "Update user_2fa request"
// @Success 200 {object} User2faResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /user_2fa/{id} [put]
func (h *Handler) UpdateUser2fa(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req User2faUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	user2fa, err := h.user2faRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "user_2fa not found")
			return
		}
		h.logger.Error("Failed to get user_2fa", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user_2fa")
		return
	}


	
	user2fa.Userid = req.Userid
	

	
	user2fa.Secret = req.Secret
	

	
	user2fa.Backupcodes = req.Backupcodes
	

	
	user2fa.Isenabled = req.Isenabled
	

	
	user2fa.Enabledat = req.Enabledat
	

	
	user2fa.Lastusedat = req.Lastusedat
	


	if err := h.user2faRepo.Update(ctx, user2fa); err != nil {
		h.logger.Error("Failed to update user_2fa", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user_2fa")
		return
	}

	response := User2faResponse{

		Id: user2fa.Id,

		Userid: user2fa.Userid,

		Secret: user2fa.Secret,

		Backupcodes: user2fa.Backupcodes,

		Isenabled: user2fa.Isenabled,

		Enabledat: user2fa.Enabledat,

		Lastusedat: user2fa.Lastusedat,

		Createdat: user2fa.Createdat,

		Updatedat: user2fa.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "user_2fa updated successfully", response)
}

// DeleteUser2fa deletes a user_2fa
// @Summary Delete user_2fa
// @Description Delete a user_2fa by ID
// @Tags user_2fa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "user_2fa ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /user_2fa/{id} [delete]
func (h *Handler) DeleteUser2fa(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.user2faRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "user_2fa not found")
			return
		}
		h.logger.Error("Failed to delete user_2fa", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user_2fa")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user_2fa deleted successfully", nil)
}
