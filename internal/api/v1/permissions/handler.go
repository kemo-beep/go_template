package permissions

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

// Handler handles permissions requests
type Handler struct {
	permissionsRepo generated.PermissionsRepository
	logger             *zap.Logger
}

// NewHandler creates a new permissions handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		permissionsRepo: generated.NewPermissionsRepository(db),
		logger:             logger,
	}
}

// CreatePermissions creates a new permissions
// @Summary Create permissions
// @Description Create a new permissions record
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PermissionsCreateRequest true "Create permissions request"
// @Success 201 {object} PermissionsResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /permissions [post]
func (h *Handler) CreatePermissions(c *gin.Context) {
	var req PermissionsCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	permissions := &generated.Permissions{

		Name: req.Name,

		Description: req.Description,

		Resource: req.Resource,

		Action: req.Action,

	}

	ctx := context.Background()
	if err := h.permissionsRepo.Create(ctx, permissions); err != nil {
		h.logger.Error("Failed to create permissions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create permissions")
		return
	}

	response := PermissionsResponse{

		Id: permissions.Id,

		Name: permissions.Name,

		Description: permissions.Description,

		Resource: permissions.Resource,

		Action: permissions.Action,

		Createdat: permissions.Createdat,

		Updatedat: permissions.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "permissions created successfully", response)
}

// GetPermissions gets a permissions by ID
// @Summary Get permissions
// @Description Get a permissions by ID
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "permissions ID"
// @Success 200 {object} PermissionsResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /permissions/{id} [get]
func (h *Handler) GetPermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	permissions, err := h.permissionsRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "permissions not found")
			return
		}
		h.logger.Error("Failed to get permissions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get permissions")
		return
	}

	response := PermissionsResponse{

		Id: permissions.Id,

		Name: permissions.Name,

		Description: permissions.Description,

		Resource: permissions.Resource,

		Action: permissions.Action,

		Createdat: permissions.Createdat,

		Updatedat: permissions.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "permissions retrieved successfully", response)
}

// GetAllPermissionss gets all permissionss with pagination
// @Summary Get all permissionss
// @Description Get all permissionss with pagination
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /permissions [get]
func (h *Handler) GetAllPermissionss(c *gin.Context) {
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
	permissionss, total, err := h.permissionsRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get permissionss", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get permissionss")
		return
	}

	var responses []PermissionsResponse
	for _, permissions := range permissionss {
		responses = append(responses, PermissionsResponse{

			Id: permissions.Id,

			Name: permissions.Name,

			Description: permissions.Description,

			Resource: permissions.Resource,

			Action: permissions.Action,

			Createdat: permissions.Createdat,

			Updatedat: permissions.Updatedat,

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

	utils.SuccessResponse(c, http.StatusOK, "permissionss retrieved successfully", pagination)
}

// UpdatePermissions updates a permissions
// @Summary Update permissions
// @Description Update a permissions by ID
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "permissions ID"
// @Param request body PermissionsUpdateRequest true "Update permissions request"
// @Success 200 {object} PermissionsResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /permissions/{id} [put]
func (h *Handler) UpdatePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req PermissionsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	permissions, err := h.permissionsRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "permissions not found")
			return
		}
		h.logger.Error("Failed to get permissions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get permissions")
		return
	}


	
	permissions.Name = req.Name
	

	
	permissions.Description = req.Description
	

	
	permissions.Resource = req.Resource
	

	
	permissions.Action = req.Action
	


	if err := h.permissionsRepo.Update(ctx, permissions); err != nil {
		h.logger.Error("Failed to update permissions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update permissions")
		return
	}

	response := PermissionsResponse{

		Id: permissions.Id,

		Name: permissions.Name,

		Description: permissions.Description,

		Resource: permissions.Resource,

		Action: permissions.Action,

		Createdat: permissions.Createdat,

		Updatedat: permissions.Updatedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "permissions updated successfully", response)
}

// DeletePermissions deletes a permissions
// @Summary Delete permissions
// @Description Delete a permissions by ID
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "permissions ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /permissions/{id} [delete]
func (h *Handler) DeletePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.permissionsRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "permissions not found")
			return
		}
		h.logger.Error("Failed to delete permissions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete permissions")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "permissions deleted successfully", nil)
}
