package roles

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

// Handler handles roles requests
type Handler struct {
	rolesRepo generated.RolesRepository
	logger             *zap.Logger
}

// NewHandler creates a new roles handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		rolesRepo: generated.NewRolesRepository(db),
		logger:             logger,
	}
}

// CreateRoles creates a new roles
// @Summary Create roles
// @Description Create a new roles record
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RolesCreateRequest true "Create roles request"
// @Success 201 {object} RolesResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /roles [post]
func (h *Handler) CreateRoles(c *gin.Context) {
	var req RolesCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	roles := &generated.Roles{

		Name: req.Name,

		Description: req.Description,

	}

	ctx := context.Background()
	if err := h.rolesRepo.Create(ctx, roles); err != nil {
		h.logger.Error("Failed to create roles", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create roles")
		return
	}

	response := RolesResponse{

		Id: roles.Id,

		Name: roles.Name,

		Description: roles.Description,

		Createdat: roles.Createdat,

		Updatedat: roles.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "roles created successfully", response)
}

// GetRoles gets a roles by ID
// @Summary Get roles
// @Description Get a roles by ID
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "roles ID"
// @Success 200 {object} RolesResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /roles/{id} [get]
func (h *Handler) GetRoles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	roles, err := h.rolesRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "roles not found")
			return
		}
		h.logger.Error("Failed to get roles", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get roles")
		return
	}

	response := RolesResponse{

		Id: roles.Id,

		Name: roles.Name,

		Description: roles.Description,

		Createdat: roles.Createdat,

		Updatedat: roles.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "roles retrieved successfully", response)
}

// GetAllRoless gets all roless with pagination
// @Summary Get all roless
// @Description Get all roless with pagination
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /roles [get]
func (h *Handler) GetAllRoless(c *gin.Context) {
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
	roless, total, err := h.rolesRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get roless", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get roless")
		return
	}

	var responses []RolesResponse
	for _, roles := range roless {
		responses = append(responses, RolesResponse{

			Id: roles.Id,

			Name: roles.Name,

			Description: roles.Description,

			Createdat: roles.Createdat,

			Updatedat: roles.Updatedat,

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

	utils.SuccessResponse(c, http.StatusOK, "roless retrieved successfully", pagination)
}

// UpdateRoles updates a roles
// @Summary Update roles
// @Description Update a roles by ID
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "roles ID"
// @Param request body RolesUpdateRequest true "Update roles request"
// @Success 200 {object} RolesResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /roles/{id} [put]
func (h *Handler) UpdateRoles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req RolesUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	roles, err := h.rolesRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "roles not found")
			return
		}
		h.logger.Error("Failed to get roles", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get roles")
		return
	}


	
	roles.Name = req.Name
	

	
	roles.Description = req.Description
	


	if err := h.rolesRepo.Update(ctx, roles); err != nil {
		h.logger.Error("Failed to update roles", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update roles")
		return
	}

	response := RolesResponse{

		Id: roles.Id,

		Name: roles.Name,

		Description: roles.Description,

		Createdat: roles.Createdat,

		Updatedat: roles.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "roles updated successfully", response)
}

// DeleteRoles deletes a roles
// @Summary Delete roles
// @Description Delete a roles by ID
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "roles ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /roles/{id} [delete]
func (h *Handler) DeleteRoles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.rolesRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "roles not found")
			return
		}
		h.logger.Error("Failed to delete roles", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete roles")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "roles deleted successfully", nil)
}
