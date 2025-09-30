package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/db/repository"
	"go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/internal/utils"
	"go-mobile-backend-template/pkg/config"
)

// Handler handles authentication requests
type Handler struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtService       *auth.JWTService
	logger           *zap.Logger
	cfg              *config.Config
}

// NewHandler creates a new auth handler
func NewHandler(db *gorm.DB, logger *zap.Logger, cfg *config.Config) *Handler {
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	return &Handler{
		userRepo:         repository.NewUserRepository(db),
		refreshTokenRepo: repository.NewRefreshTokenRepository(db),
		jwtService:       jwtService,
		logger:           logger,
		cfg:              cfg,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration request"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if user already exists
	ctx := context.Background()
	existingUser, _ := h.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		utils.ErrorResponse(c, http.StatusConflict, "User already exists")
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		h.logger.Error("Failed to hash password", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Create user
	user := &repository.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		IsActive: true,
		IsAdmin:  false,
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate tokens
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		h.logger.Error("Failed to generate access token", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		h.logger.Error("Failed to generate refresh token", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Save refresh token
	refreshTokenRecord := &repository.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: h.jwtService.GetRefreshTokenExpiration(),
		IsRevoked: false,
	}
	if err := h.refreshTokenRepo.Create(ctx, refreshTokenRecord); err != nil {
		h.logger.Error("Failed to save refresh token", zap.Error(err))
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserData{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			IsActive: user.IsActive,
			IsAdmin:  user.IsAdmin,
		},
		ExpiresIn: h.cfg.JWT.AccessTokenExpireInt * 60,
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	user, err := h.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check password
	if err := auth.CheckPassword(user.Password, req.Password); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check if user is active
	if !user.IsActive {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User account is inactive")
		return
	}

	// Generate tokens
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		h.logger.Error("Failed to generate access token", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		h.logger.Error("Failed to generate refresh token", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Save refresh token
	refreshTokenRecord := &repository.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: h.jwtService.GetRefreshTokenExpiration(),
		IsRevoked: false,
	}
	if err := h.refreshTokenRepo.Create(ctx, refreshTokenRecord); err != nil {
		h.logger.Error("Failed to save refresh token", zap.Error(err))
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserData{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			IsActive: user.IsActive,
			IsAdmin:  user.IsAdmin,
		},
		ExpiresIn: h.cfg.JWT.AccessTokenExpireInt * 60,
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Check if refresh token exists and is not revoked
	ctx := context.Background()
	refreshTokenRecord, err := h.refreshTokenRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Validate that the refresh token belongs to the user from JWT claims
	if refreshTokenRecord.UserID != claims.UserID {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Get user
	user, err := h.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	// Generate new access token
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		h.logger.Error("Failed to generate access token", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Optionally generate new refresh token
	newRefreshToken, err := h.jwtService.GenerateRefreshToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		h.logger.Error("Failed to generate refresh token", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Revoke old refresh token
	if err := h.refreshTokenRepo.Revoke(ctx, req.RefreshToken); err != nil {
		h.logger.Error("Failed to revoke refresh token", zap.Error(err))
	}

	// Save new refresh token
	newRefreshTokenRecord := &repository.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: h.jwtService.GetRefreshTokenExpiration(),
		IsRevoked: false,
	}
	if err := h.refreshTokenRepo.Create(ctx, newRefreshTokenRecord); err != nil {
		h.logger.Error("Failed to save refresh token", zap.Error(err))
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: UserData{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			IsActive: user.IsActive,
			IsAdmin:  user.IsAdmin,
		},
		ExpiresIn: h.cfg.JWT.AccessTokenExpireInt * 60,
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", response)
}

// Logout handles user logout
// @Summary User logout
// @Description Revoke user's refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx := context.Background()
	if err := h.refreshTokenRepo.RevokeAllForUser(ctx, userID.(uint)); err != nil {
		h.logger.Error("Failed to revoke refresh tokens", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to logout")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}
