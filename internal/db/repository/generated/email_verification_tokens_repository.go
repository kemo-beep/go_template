package generated

import (
	"context"
	"gorm.io/gorm"
)

// EmailverificationtokensRepository interface for email_verification_tokens operations
type EmailverificationtokensRepository interface {
	Create(ctx context.Context, emailVerificationTokens *Emailverificationtokens) error
	GetByID(ctx context.Context, id uint) (*Emailverificationtokens, error)
	GetAll(ctx context.Context, limit, offset int) ([]Emailverificationtokens, int64, error)
	Update(ctx context.Context, emailVerificationTokens *Emailverificationtokens) error
	Delete(ctx context.Context, id uint) error
}

// emailVerificationTokensRepository implements EmailverificationtokensRepository
type emailVerificationTokensRepository struct {
	db *gorm.DB
}

// NewEmailverificationtokensRepository creates a new EmailverificationtokensRepository
func NewEmailverificationtokensRepository(db *gorm.DB) EmailverificationtokensRepository {
	return &emailVerificationTokensRepository{db: db}
}

// Create creates a new emailVerificationTokens
func (r *emailVerificationTokensRepository) Create(ctx context.Context, emailVerificationTokens *Emailverificationtokens) error {
	return r.db.WithContext(ctx).Create(emailVerificationTokens).Error
}

// GetByID gets a emailVerificationTokens by ID
func (r *emailVerificationTokensRepository) GetByID(ctx context.Context, id uint) (*Emailverificationtokens, error) {
	var emailVerificationTokens Emailverificationtokens
	err := r.db.WithContext(ctx).First(&emailVerificationTokens, id).Error
	if err != nil {
		return nil, err
	}
	return &emailVerificationTokens, nil
}

// GetAll gets all emailVerificationTokenss with pagination
func (r *emailVerificationTokensRepository) GetAll(ctx context.Context, limit, offset int) ([]Emailverificationtokens, int64, error) {
	var emailVerificationTokenss []Emailverificationtokens
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Emailverificationtokens{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&emailVerificationTokenss).Error
	return emailVerificationTokenss, total, err
}

// Update updates a emailVerificationTokens
func (r *emailVerificationTokensRepository) Update(ctx context.Context, emailVerificationTokens *Emailverificationtokens) error {
	return r.db.WithContext(ctx).Save(emailVerificationTokens).Error
}

// Delete deletes a emailVerificationTokens by ID
func (r *emailVerificationTokensRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Emailverificationtokens{}, id).Error
}
