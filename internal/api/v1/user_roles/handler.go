package user_roles

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

// Handler handles user_roles requests
type Handler struct {
	userRolesRepo generated.UserrolesRepository
	logger             *zap.Logger
}

// NewHandler creates a new user_roles handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		userRolesRepo: generated.NewUserrolesRepository(db),
		logger:             logger,
	}
}

// CreateUserroles creates a new user_roles
// @Summary Create user_roles
// @Description Create a new user_roles record
// @Tags user_roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UserrolesCreateRequest true "Create user_roles request"
// @Success 201 {object} UserrolesResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /user_roles [post]
func (h *Handler) CreateUserroles(c *gin.Context) {
	var req UserrolesCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userRoles := &generated.Userroles{

		Userid: req.Userid,

		Roleid: req.Roleid,

		Assignedat: req.Assignedat,

		Assignedby: req.Assignedby,

	}

	ctx := context.Background()
	if err := h.userRolesRepo.Create(ctx, userRoles); err != nil {
		h.logger.Error("Failed to create user_roles", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user_roles")
		return
	}

	response := UserrolesResponse{

		Userid: userRoles.Userid,

		Roleid: userRoles.Roleid,

		Assignedat: userRoles.Assignedat,

		Assignedby: userRoles.Assignedby,

	}

	utils.SuccessResponse(c, http.StatusCreated, "user_roles created successfully", response)
}

// GetUserroles gets a user_roles by ID
// @Summary Get user_roles
// @Description Get a user_roles by ID
// @Tags user_roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "user_roles ID"
// @Success 200 {object} UserrolesResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /user_roles/{id} [get]
func (h *Handler) GetUserroles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	userRoles, err := h.userRolesRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "user_roles not found")
			return
		}
		h.logger.Error("Failed to get user_roles", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user_roles")
		return
	}

	response := UserrolesResponse{

		Userid: userRoles.Userid,

		Roleid: userRoles.Roleid,

		Assignedat: userRoles.Assignedat,

		Assignedby: userRoles.Assignedby,

	}

	utils.SuccessResponse(c, http.StatusOK, "user_roles retrieved successfully", response)
}

// GetAllUserroless gets all user_roless with pagination
// @Summary Get all user_roless
// @Description Get all user_roless with pagination
// @Tags user_roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /user_roles [get]
func (h *Handler) GetAllUserroless(c *gin.Context) {
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
	userRoless, total, err := h.userRolesRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get user_roless", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user_roless")
		return
	}

	var responses []UserrolesResponse
	for _, userRoles := range userRoless {
		responses = append(responses, UserrolesResponse{

			Userid: userRoles.Userid,

			Roleid: userRoles.Roleid,

			Assignedat: userRoles.Assignedat,

			Assignedby: userRoles.Assignedby,

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

	utils.SuccessResponse(c, http.StatusOK, "user_roless retrieved successfully", pagination)
}

// UpdateUserroles updates a user_roles
// @Summary Update user_roles
// @Description Update a user_roles by ID
// @Tags user_roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "user_roles ID"
// @Param request body UserrolesUpdateRequest true "Update user_roles request"
// @Success 200 {object} UserrolesResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /user_roles/{id} [put]
func (h *Handler) UpdateUserroles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req UserrolesUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	userRoles, err := h.userRolesRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "user_roles not found")
			return
		}
		h.logger.Error("Failed to get user_roles", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user_roles")
		return
	}


	
	userRoles.Userid = req.Userid
	

	
	userRoles.Roleid = req.Roleid
	

	
	userRoles.Assignedat = req.Assignedat
	

	
	userRoles.Assignedby = req.Assignedby
	


	if err := h.userRolesRepo.Update(ctx, userRoles); err != nil {
		h.logger.Error("Failed to update user_roles", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user_roles")
		return
	}

	response := UserrolesResponse{

		Userid: userRoles.Userid,

		Roleid: userRoles.Roleid,

		Assignedat: userRoles.Assignedat,

		Assignedby: userRoles.Assignedby,

	}

	utils.SuccessResponse(c, http.StatusOK, "user_roles updated successfully", response)
}

// DeleteUserroles deletes a user_roles
// @Summary Delete user_roles
// @Description Delete a user_roles by ID
// @Tags user_roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "user_roles ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /user_roles/{id} [delete]
func (h *Handler) DeleteUserroles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.userRolesRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "user_roles not found")
			return
		}
		h.logger.Error("Failed to delete user_roles", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user_roles")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user_roles deleted successfully", nil)
}
