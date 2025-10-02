package role_permissions

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

// Handler handles role_permissions requests
type Handler struct {
	rolePermissionsRepo generated.RolepermissionsRepository
	logger             *zap.Logger
}

// NewHandler creates a new role_permissions handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		rolePermissionsRepo: generated.NewRolepermissionsRepository(db),
		logger:             logger,
	}
}

// CreateRolepermissions creates a new role_permissions
// @Summary Create role_permissions
// @Description Create a new role_permissions record
// @Tags role_permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RolepermissionsCreateRequest true "Create role_permissions request"
// @Success 201 {object} RolepermissionsResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /role_permissions [post]
func (h *Handler) CreateRolepermissions(c *gin.Context) {
	var req RolepermissionsCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	rolePermissions := &generated.Rolepermissions{

		Roleid: req.Roleid,

		Permissionid: req.Permissionid,

	}

	ctx := context.Background()
	if err := h.rolePermissionsRepo.Create(ctx, rolePermissions); err != nil {
		h.logger.Error("Failed to create role_permissions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create role_permissions")
		return
	}

	response := RolepermissionsResponse{

		Roleid: rolePermissions.Roleid,

		Permissionid: rolePermissions.Permissionid,

		Createdat: rolePermissions.Createdat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "role_permissions created successfully", response)
}

// GetRolepermissions gets a role_permissions by ID
// @Summary Get role_permissions
// @Description Get a role_permissions by ID
// @Tags role_permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "role_permissions ID"
// @Success 200 {object} RolepermissionsResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /role_permissions/{id} [get]
func (h *Handler) GetRolepermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	rolePermissions, err := h.rolePermissionsRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "role_permissions not found")
			return
		}
		h.logger.Error("Failed to get role_permissions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get role_permissions")
		return
	}

	response := RolepermissionsResponse{

		Roleid: rolePermissions.Roleid,

		Permissionid: rolePermissions.Permissionid,

		Createdat: rolePermissions.Createdat,

	}

	utils.SuccessResponse(c, http.StatusOK, "role_permissions retrieved successfully", response)
}

// GetAllRolepermissionss gets all role_permissionss with pagination
// @Summary Get all role_permissionss
// @Description Get all role_permissionss with pagination
// @Tags role_permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /role_permissions [get]
func (h *Handler) GetAllRolepermissionss(c *gin.Context) {
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
	rolePermissionss, total, err := h.rolePermissionsRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get role_permissionss", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get role_permissionss")
		return
	}

	var responses []RolepermissionsResponse
	for _, rolePermissions := range rolePermissionss {
		responses = append(responses, RolepermissionsResponse{

			Roleid: rolePermissions.Roleid,

			Permissionid: rolePermissions.Permissionid,

			Createdat: rolePermissions.Createdat,

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

	utils.SuccessResponse(c, http.StatusOK, "role_permissionss retrieved successfully", pagination)
}

// UpdateRolepermissions updates a role_permissions
// @Summary Update role_permissions
// @Description Update a role_permissions by ID
// @Tags role_permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "role_permissions ID"
// @Param request body RolepermissionsUpdateRequest true "Update role_permissions request"
// @Success 200 {object} RolepermissionsResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /role_permissions/{id} [put]
func (h *Handler) UpdateRolepermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req RolepermissionsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	rolePermissions, err := h.rolePermissionsRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "role_permissions not found")
			return
		}
		h.logger.Error("Failed to get role_permissions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get role_permissions")
		return
	}


	
	rolePermissions.Roleid = req.Roleid
	

	
	rolePermissions.Permissionid = req.Permissionid
	


	if err := h.rolePermissionsRepo.Update(ctx, rolePermissions); err != nil {
		h.logger.Error("Failed to update role_permissions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update role_permissions")
		return
	}

	response := RolepermissionsResponse{

		Roleid: rolePermissions.Roleid,

		Permissionid: rolePermissions.Permissionid,

		Createdat: rolePermissions.Createdat,

	}

	utils.SuccessResponse(c, http.StatusOK, "role_permissions updated successfully", response)
}

// DeleteRolepermissions deletes a role_permissions
// @Summary Delete role_permissions
// @Description Delete a role_permissions by ID
// @Tags role_permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "role_permissions ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /role_permissions/{id} [delete]
func (h *Handler) DeleteRolepermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.rolePermissionsRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "role_permissions not found")
			return
		}
		h.logger.Error("Failed to delete role_permissions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete role_permissions")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "role_permissions deleted successfully", nil)
}
