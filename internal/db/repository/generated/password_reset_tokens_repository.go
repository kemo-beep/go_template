package generated

import (
	"context"
	"gorm.io/gorm"
)

// PasswordresettokensRepository interface for password_reset_tokens operations
type PasswordresettokensRepository interface {
	Create(ctx context.Context, passwordResetTokens *Passwordresettokens) error
	GetByID(ctx context.Context, id uint) (*Passwordresettokens, error)
	GetAll(ctx context.Context, limit, offset int) ([]Passwordresettokens, int64, error)
	Update(ctx context.Context, passwordResetTokens *Passwordresettokens) error
	Delete(ctx context.Context, id uint) error
}

// passwordResetTokensRepository implements PasswordresettokensRepository
type passwordResetTokensRepository struct {
	db *gorm.DB
}

// NewPasswordresettokensRepository creates a new PasswordresettokensRepository
func NewPasswordresettokensRepository(db *gorm.DB) PasswordresettokensRepository {
	return &passwordResetTokensRepository{db: db}
}

// Create creates a new passwordResetTokens
func (r *passwordResetTokensRepository) Create(ctx context.Context, passwordResetTokens *Passwordresettokens) error {
	return r.db.WithContext(ctx).Create(passwordResetTokens).Error
}

// GetByID gets a passwordResetTokens by ID
func (r *passwordResetTokensRepository) GetByID(ctx context.Context, id uint) (*Passwordresettokens, error) {
	var passwordResetTokens Passwordresettokens
	err := r.db.WithContext(ctx).First(&passwordResetTokens, id).Error
	if err != nil {
		return nil, err
	}
	return &passwordResetTokens, nil
}

// GetAll gets all passwordResetTokenss with pagination
func (r *passwordResetTokensRepository) GetAll(ctx context.Context, limit, offset int) ([]Passwordresettokens, int64, error) {
	var passwordResetTokenss []Passwordresettokens
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Passwordresettokens{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&passwordResetTokenss).Error
	return passwordResetTokenss, total, err
}

// Update updates a passwordResetTokens
func (r *passwordResetTokensRepository) Update(ctx context.Context, passwordResetTokens *Passwordresettokens) error {
	return r.db.WithContext(ctx).Save(passwordResetTokens).Error
}

// Delete deletes a passwordResetTokens by ID
func (r *passwordResetTokensRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Passwordresettokens{}, id).Error
}
