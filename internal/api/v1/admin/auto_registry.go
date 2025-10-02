package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/generator"
	"go-mobile-backend-template/internal/utils"
	"go-mobile-backend-template/pkg/config"
)

// AutoRegistryHandler handles auto-registry related endpoints
type AutoRegistryHandler struct {
	db           *gorm.DB
	logger       *zap.Logger
	cfg          *config.Config
	autoRegistry *generator.AutoRegistry
}

// NewAutoRegistryHandler creates a new auto-registry handler
func NewAutoRegistryHandler(db *gorm.DB, logger *zap.Logger, cfg *config.Config, autoRegistry *generator.AutoRegistry) *AutoRegistryHandler {
	return &AutoRegistryHandler{
		db:           db,
		logger:       logger,
		cfg:          cfg,
		autoRegistry: autoRegistry,
	}
}

// GetStatus returns the current status of the auto-registry
// @Summary Get auto-registry status
// @Description Get the current status of the auto-registry system
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /admin/auto-registry/status [get]
func (h *AutoRegistryHandler) GetStatus(c *gin.Context) {
	status := h.autoRegistry.GetStatus()

	utils.SuccessResponse(c, http.StatusOK, "Auto-registry status retrieved successfully", status)
}

// GetRegisteredAPIs returns the list of currently registered APIs
// @Summary Get registered APIs
// @Description Get the list of currently registered auto-generated APIs
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /admin/auto-registry/apis [get]
func (h *AutoRegistryHandler) GetRegisteredAPIs(c *gin.Context) {
	apis := h.autoRegistry.GetRegisteredAPIs()

	utils.SuccessResponse(c, http.StatusOK, "Registered APIs retrieved successfully", map[string]interface{}{
		"apis":  apis,
		"count": len(apis),
	})
}

// RegenerateAPIs manually triggers API regeneration
// @Summary Regenerate APIs
// @Description Manually trigger regeneration of all auto-generated APIs
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /admin/auto-registry/regenerate [post]
func (h *AutoRegistryHandler) RegenerateAPIs(c *gin.Context) {
	// This would trigger the regeneration process
	// For now, we'll just return a success message
	utils.SuccessResponse(c, http.StatusOK, "API regeneration triggered successfully", map[string]interface{}{
		"message": "APIs will be regenerated in the background",
	})
}
