package admin

import (
	"net/http"
	"strconv"

	"go-mobile-backend-template/internal/db/repository"
	"go-mobile-backend-template/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRoleHandler(db *gorm.DB, logger *zap.Logger) *RoleHandler {
	return &RoleHandler{
		db:     db,
		logger: logger,
	}
}

// ListRoles godoc
// @Summary List all roles (Admin)
// @Description Get list of all roles with permissions
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /admin/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	roleRepo := repository.NewRoleRepository(h.db)

	roles, total, err := roleRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list roles", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to fetch roles"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Roles fetched successfully", gin.H{
		"roles": roles,
		"total": total,
		"page":  page,
		"limit": limit,
	}))
}

// CreateRoleRequest represents role creation request
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateRole godoc
// @Summary Create new role (Admin)
// @Description Create a new role
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRoleRequest true "Role creation request"
// @Success 201 {object} utils.Response
// @Router /admin/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request"))
		return
	}

	role := &repository.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	roleRepo := repository.NewRoleRepository(h.db)

	if err := roleRepo.Create(c.Request.Context(), role); err != nil {
		h.logger.Error("Failed to create role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to create role"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponseData("Role created successfully", role))
}

// GetRole godoc
// @Summary Get role details (Admin)
// @Description Get role with permissions
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} utils.Response
// @Router /admin/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid role ID"))
		return
	}

	roleRepo := repository.NewRoleRepository(h.db)

	role, err := roleRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponseData("Role not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Role fetched successfully", role))
}

// AssignPermissionsRequest represents permission assignment request
type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

// AssignPermissions godoc
// @Summary Assign permissions to role (Admin)
// @Description Assign multiple permissions to a role
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param request body AssignPermissionsRequest true "Permissions request"
// @Success 200 {object} utils.Response
// @Router /admin/roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid role ID"))
		return
	}

	var req AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request"))
		return
	}

	roleRepo := repository.NewRoleRepository(h.db)

	if err := roleRepo.AssignPermissions(c.Request.Context(), uint(id), req.PermissionIDs); err != nil {
		h.logger.Error("Failed to assign permissions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to assign permissions"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Permissions assigned successfully", nil))
}

// ListPermissions godoc
// @Summary List all permissions (Admin)
// @Description Get list of all available permissions
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /admin/permissions [get]
func (h *RoleHandler) ListPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset := (page - 1) * limit

	permRepo := repository.NewPermissionRepository(h.db)

	permissions, total, err := permRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list permissions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to fetch permissions"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Permissions fetched successfully", gin.H{
		"permissions": permissions,
		"total":       total,
	}))
}
