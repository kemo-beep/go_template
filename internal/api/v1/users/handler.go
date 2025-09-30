package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/db/repository"
	"go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/internal/utils"
)

// Handler handles user requests
type Handler struct {
	userRepo repository.UserRepository
	logger   *zap.Logger
}

// NewHandler creates a new user handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		userRepo: repository.NewUserRepository(db),
		logger:   logger,
	}
}

// GetProfile gets the current user's profile
// @Summary Get user profile
// @Description Get current authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} utils.Response
// @Router /users/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx := context.Background()
	user, err := h.userRepo.GetByID(ctx, userID.(uint))
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		IsActive: user.IsActive,
		IsAdmin:  user.IsAdmin,
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", response)
}

// UpdateProfile updates the current user's profile
// @Summary Update user profile
// @Description Update current authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Update profile request"
// @Success 200 {object} UserResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /users/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	user, err := h.userRepo.GetByID(ctx, userID.(uint))
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}

	if err := h.userRepo.Update(ctx, user); err != nil {
		h.logger.Error("Failed to update user", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		IsActive: user.IsActive,
		IsAdmin:  user.IsAdmin,
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", response)
}

// ChangePassword changes the current user's password
// @Summary Change password
// @Description Change current authenticated user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Change password request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /users/me/change-password [post]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	user, err := h.userRepo.GetByID(ctx, userID.(uint))
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Verify old password
	if err := auth.CheckPassword(user.Password, req.OldPassword); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid old password")
		return
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		h.logger.Error("Failed to hash password", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to change password")
		return
	}

	user.Password = hashedPassword
	if err := h.userRepo.Update(ctx, user); err != nil {
		h.logger.Error("Failed to update password", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to change password")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

// DeleteAccount deletes the current user's account
// @Summary Delete account
// @Description Delete current authenticated user's account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /users/me [delete]
func (h *Handler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx := context.Background()
	if err := h.userRepo.Delete(ctx, userID.(uint)); err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account deleted successfully", nil)
}
