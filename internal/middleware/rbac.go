package middleware

import (
	"net/http"
	"strings"

	"go-mobile-backend-template/internal/db/repository"
	"go-mobile-backend-template/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RequirePermission middleware checks if user has required permission
func RequirePermission(resource, action string, db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			logger.Warn("User ID not found in context")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			logger.Warn("Invalid user ID type in context")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// Check if user is admin (admins have all permissions)
		isAdmin, exists := c.Get("is_admin")
		if exists && isAdmin.(bool) {
			c.Next()
			return
		}

		// Check permission
		permRepo := repository.NewPermissionRepository(db)
		hasPermission, err := permRepo.CheckUserPermission(c.Request.Context(), uid, resource, action)
		if err != nil {
			logger.Error("Failed to check permission", zap.Error(err))
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check permissions")
			c.Abort()
			return
		}

		if !hasPermission {
			logger.Warn("User lacks required permission",
				zap.Uint("user_id", uid),
				zap.String("resource", resource),
				zap.String("action", action),
			)
			utils.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(roleName string, db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			logger.Warn("User ID not found in context")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			logger.Warn("Invalid user ID type in context")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// Check if user is admin
		if roleName != "admin" {
			isAdmin, exists := c.Get("is_admin")
			if exists && isAdmin.(bool) {
				c.Next()
				return
			}
		}

		// Check role
		roleRepo := repository.NewRoleRepository(db)
		roles, err := roleRepo.GetUserRoles(c.Request.Context(), uid)
		if err != nil {
			logger.Error("Failed to get user roles", zap.Error(err))
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check roles")
			c.Abort()
			return
		}

		// Check if user has the required role
		hasRole := false
		for _, role := range roles {
			if strings.EqualFold(role.Name, roleName) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			logger.Warn("User lacks required role",
				zap.Uint("user_id", uid),
				zap.String("required_role", roleName),
			)
			utils.ErrorResponse(c, http.StatusForbidden, "Insufficient role privileges")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware ensures user is an admin
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
