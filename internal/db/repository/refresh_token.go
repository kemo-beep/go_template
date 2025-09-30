package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// refreshTokenRepository implements RefreshTokenRepository interface
type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *refreshTokenRepository) Create(ctx context.Context, token *RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

// GetByToken retrieves a refresh token by token string
func (r *refreshTokenRepository) GetByToken(ctx context.Context, token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	if err := r.db.WithContext(ctx).
		Where("token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&refreshToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	return &refreshToken, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (r *refreshTokenRepository) GetByUserID(ctx context.Context, userID uint) ([]*RefreshToken, error) {
	var tokens []*RefreshToken
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to get refresh tokens by user ID: %w", err)
	}
	return tokens, nil
}

// Revoke revokes a specific refresh token
func (r *refreshTokenRepository) Revoke(ctx context.Context, token string) error {
	if err := r.db.WithContext(ctx).
		Model(&RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error; err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *refreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).
		Model(&RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error; err != nil {
		return fmt.Errorf("failed to revoke all refresh tokens for user: %w", err)
	}
	return nil
}

// DeleteExpired deletes expired refresh tokens
func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	if err := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}
	return nil
}
