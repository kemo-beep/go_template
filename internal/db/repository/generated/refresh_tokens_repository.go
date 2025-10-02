package generated

import (
	"context"
	"gorm.io/gorm"
)

// RefreshtokensRepository interface for refresh_tokens operations
type RefreshtokensRepository interface {
	Create(ctx context.Context, refreshTokens *Refreshtokens) error
	GetByID(ctx context.Context, id uint) (*Refreshtokens, error)
	GetAll(ctx context.Context, limit, offset int) ([]Refreshtokens, int64, error)
	Update(ctx context.Context, refreshTokens *Refreshtokens) error
	Delete(ctx context.Context, id uint) error
}

// refreshTokensRepository implements RefreshtokensRepository
type refreshTokensRepository struct {
	db *gorm.DB
}

// NewRefreshtokensRepository creates a new RefreshtokensRepository
func NewRefreshtokensRepository(db *gorm.DB) RefreshtokensRepository {
	return &refreshTokensRepository{db: db}
}

// Create creates a new refreshTokens
func (r *refreshTokensRepository) Create(ctx context.Context, refreshTokens *Refreshtokens) error {
	return r.db.WithContext(ctx).Create(refreshTokens).Error
}

// GetByID gets a refreshTokens by ID
func (r *refreshTokensRepository) GetByID(ctx context.Context, id uint) (*Refreshtokens, error) {
	var refreshTokens Refreshtokens
	err := r.db.WithContext(ctx).First(&refreshTokens, id).Error
	if err != nil {
		return nil, err
	}
	return &refreshTokens, nil
}

// GetAll gets all refreshTokenss with pagination
func (r *refreshTokensRepository) GetAll(ctx context.Context, limit, offset int) ([]Refreshtokens, int64, error) {
	var refreshTokenss []Refreshtokens
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Refreshtokens{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&refreshTokenss).Error
	return refreshTokenss, total, err
}

// Update updates a refreshTokens
func (r *refreshTokensRepository) Update(ctx context.Context, refreshTokens *Refreshtokens) error {
	return r.db.WithContext(ctx).Save(refreshTokens).Error
}

// Delete deletes a refreshTokens by ID
func (r *refreshTokensRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Refreshtokens{}, id).Error
}
