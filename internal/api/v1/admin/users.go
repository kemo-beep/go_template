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

type UserHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserHandler(db *gorm.DB, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		db:     db,
		logger: logger,
	}
}

// ListUsers godoc
// @Summary List all users (Admin)
// @Description Get paginated list of all users with their roles
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param search query string false "Search by email or name"
// @Param role query string false "Filter by role name"
// @Param is_active query boolean false "Filter by active status"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	roleFilter := c.Query("role")
	isActiveStr := c.Query("is_active")

	offset := (page - 1) * limit
	if limit > 100 {
		limit = 100
	}

	roleRepo := repository.NewRoleRepository(h.db)

	// Build query
	query := h.db.Model(&repository.User{})

	// Apply filters
	if search != "" {
		query = query.Where("email LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if isActiveStr != "" {
		isActive, _ := strconv.ParseBool(isActiveStr)
		query = query.Where("is_active = ?", isActive)
	}

	// Get users
	var users []repository.User
	var total int64

	if err := query.Count(&total).Error; err != nil {
		h.logger.Error("Failed to count users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to fetch users"))
		return
	}

	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		h.logger.Error("Failed to fetch users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to fetch users"))
		return
	}

	// Get roles for each user
	type UserWithRoles struct {
		repository.User
		Roles []string `json:"roles"`
	}

	usersWithRoles := make([]UserWithRoles, len(users))
	for i, user := range users {
		roles, err := roleRepo.GetUserRoles(c.Request.Context(), user.ID)
		if err != nil {
			h.logger.Error("Failed to get user roles", zap.Error(err), zap.Uint("user_id", user.ID))
			continue
		}

		roleNames := make([]string, len(roles))
		for j, role := range roles {
			roleNames[j] = role.Name
		}

		usersWithRoles[i] = UserWithRoles{
			User:  user,
			Roles: roleNames,
		}
	}

	// Filter by role if specified
	if roleFilter != "" {
		filtered := []UserWithRoles{}
		for _, u := range usersWithRoles {
			for _, r := range u.Roles {
				if r == roleFilter {
					filtered = append(filtered, u)
					break
				}
			}
		}
		usersWithRoles = filtered
		total = int64(len(filtered))
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Users fetched successfully", gin.H{
		"users": usersWithRoles,
		"total": total,
		"page":  page,
		"limit": limit,
	}))
}

// GetUser godoc
// @Summary Get user details (Admin)
// @Description Get detailed information about a specific user
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid user ID"))
		return
	}

	userRepo := repository.NewUserRepository(h.db)
	roleRepo := repository.NewRoleRepository(h.db)

	user, err := userRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		c.JSON(http.StatusNotFound, utils.ErrorResponseData("User not found"))
		return
	}

	// Get user roles
	roles, err := roleRepo.GetUserRoles(c.Request.Context(), user.ID)
	if err != nil {
		h.logger.Error("Failed to get user roles", zap.Error(err))
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("User fetched successfully", gin.H{
		"user":  user,
		"roles": roles,
	}))
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
	IsAdmin  *bool   `json:"is_admin"`
}

// UpdateUser godoc
// @Summary Update user (Admin)
// @Description Update user information
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "Update request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid user ID"))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request"))
		return
	}

	userRepo := repository.NewUserRepository(h.db)

	user, err := userRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponseData("User not found"))
		return
	}

	// Update fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}

	if err := userRepo.Update(c.Request.Context(), user); err != nil {
		h.logger.Error("Failed to update user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to update user"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("User updated successfully", user))
}

// AssignRoleRequest represents role assignment request
type AssignRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

// AssignRole godoc
// @Summary Assign role to user (Admin)
// @Description Assign a role to a user
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body AssignRoleRequest true "Role assignment request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /admin/users/{id}/roles [post]
func (h *UserHandler) AssignRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid user ID"))
		return
	}

	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request"))
		return
	}

	roleRepo := repository.NewRoleRepository(h.db)

	// Get current admin user ID
	adminUserID := c.GetUint("user_id")

	if err := roleRepo.AssignRoleToUser(c.Request.Context(), uint(userID), req.RoleID, &adminUserID); err != nil {
		h.logger.Error("Failed to assign role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to assign role"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Role assigned successfully", nil))
}

// RemoveRole godoc
// @Summary Remove role from user (Admin)
// @Description Remove a role from a user
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param roleId path int true "Role ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /admin/users/{id}/roles/{roleId} [delete]
func (h *UserHandler) RemoveRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid user ID"))
		return
	}

	roleID, err := strconv.ParseUint(c.Param("roleId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid role ID"))
		return
	}

	roleRepo := repository.NewRoleRepository(h.db)

	if err := roleRepo.RemoveRoleFromUser(c.Request.Context(), uint(userID), uint(roleID)); err != nil {
		h.logger.Error("Failed to remove role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to remove role"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Role removed successfully", nil))
}

// DeleteUser godoc
// @Summary Delete user (Admin)
// @Description Soft delete a user
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid user ID"))
		return
	}

	// Prevent self-deletion
	currentUserID := c.GetUint("user_id")
	if currentUserID == uint(id) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Cannot delete your own account"))
		return
	}

	userRepo := repository.NewUserRepository(h.db)

	if err := userRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to delete user"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("User deleted successfully", nil))
}
