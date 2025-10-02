package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/pkg/config"

	"go-mobile-backend-template/internal/api/v1/role_permissions"

	"go-mobile-backend-template/internal/api/v1/test_table"

	"go-mobile-backend-template/internal/api/v1/users"

	"go-mobile-backend-template/internal/api/v1/api_keys"

	"go-mobile-backend-template/internal/api/v1/email_verification_tokens"

	"go-mobile-backend-template/internal/api/v1/oauth_providers"

	"go-mobile-backend-template/internal/api/v1/permissions"

	"go-mobile-backend-template/internal/api/v1/roles"

	"go-mobile-backend-template/internal/api/v1/sessions"

	"go-mobile-backend-template/internal/api/v1/user_2fa"

	"go-mobile-backend-template/internal/api/v1/user_roles"

	"go-mobile-backend-template/internal/api/v1/files"

	"go-mobile-backend-template/internal/api/v1/password_reset_tokens"

	"go-mobile-backend-template/internal/api/v1/refresh_tokens"

)

// RegisterGeneratedRoutes registers all auto-generated API routes
func RegisterGeneratedRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {

	role_permissions.RegisterRoutes(router, db, logger, cfg)

	test_table.RegisterRoutes(router, db, logger, cfg)

	users.RegisterRoutes(router, db, logger, cfg)

	api_keys.RegisterRoutes(router, db, logger, cfg)

	email_verification_tokens.RegisterRoutes(router, db, logger, cfg)

	oauth_providers.RegisterRoutes(router, db, logger, cfg)

	permissions.RegisterRoutes(router, db, logger, cfg)

	roles.RegisterRoutes(router, db, logger, cfg)

	sessions.RegisterRoutes(router, db, logger, cfg)

	user_2fa.RegisterRoutes(router, db, logger, cfg)

	user_roles.RegisterRoutes(router, db, logger, cfg)

	files.RegisterRoutes(router, db, logger, cfg)

	password_reset_tokens.RegisterRoutes(router, db, logger, cfg)

	refresh_tokens.RegisterRoutes(router, db, logger, cfg)

}
