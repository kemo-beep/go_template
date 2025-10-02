package generated

import (
	"context"
	"gorm.io/gorm"
)

// ApikeysRepository interface for api_keys operations
type ApikeysRepository interface {
	Create(ctx context.Context, apiKeys *Apikeys) error
	GetByID(ctx context.Context, id uint) (*Apikeys, error)
	GetAll(ctx context.Context, limit, offset int) ([]Apikeys, int64, error)
	Update(ctx context.Context, apiKeys *Apikeys) error
	Delete(ctx context.Context, id uint) error
}

// apiKeysRepository implements ApikeysRepository
type apiKeysRepository struct {
	db *gorm.DB
}

// NewApikeysRepository creates a new ApikeysRepository
func NewApikeysRepository(db *gorm.DB) ApikeysRepository {
	return &apiKeysRepository{db: db}
}

// Create creates a new apiKeys
func (r *apiKeysRepository) Create(ctx context.Context, apiKeys *Apikeys) error {
	return r.db.WithContext(ctx).Create(apiKeys).Error
}

// GetByID gets a apiKeys by ID
func (r *apiKeysRepository) GetByID(ctx context.Context, id uint) (*Apikeys, error) {
	var apiKeys Apikeys
	err := r.db.WithContext(ctx).First(&apiKeys, id).Error
	if err != nil {
		return nil, err
	}
	return &apiKeys, nil
}

// GetAll gets all apiKeyss with pagination
func (r *apiKeysRepository) GetAll(ctx context.Context, limit, offset int) ([]Apikeys, int64, error) {
	var apiKeyss []Apikeys
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Apikeys{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&apiKeyss).Error
	return apiKeyss, total, err
}

// Update updates a apiKeys
func (r *apiKeysRepository) Update(ctx context.Context, apiKeys *Apikeys) error {
	return r.db.WithContext(ctx).Save(apiKeys).Error
}

// Delete deletes a apiKeys by ID
func (r *apiKeysRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Apikeys{}, id).Error
}
