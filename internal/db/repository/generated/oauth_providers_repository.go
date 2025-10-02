package generated

import (
	"context"
	"gorm.io/gorm"
)

// OauthprovidersRepository interface for oauth_providers operations
type OauthprovidersRepository interface {
	Create(ctx context.Context, oauthProviders *Oauthproviders) error
	GetByID(ctx context.Context, id uint) (*Oauthproviders, error)
	GetAll(ctx context.Context, limit, offset int) ([]Oauthproviders, int64, error)
	Update(ctx context.Context, oauthProviders *Oauthproviders) error
	Delete(ctx context.Context, id uint) error
}

// oauthProvidersRepository implements OauthprovidersRepository
type oauthProvidersRepository struct {
	db *gorm.DB
}

// NewOauthprovidersRepository creates a new OauthprovidersRepository
func NewOauthprovidersRepository(db *gorm.DB) OauthprovidersRepository {
	return &oauthProvidersRepository{db: db}
}

// Create creates a new oauthProviders
func (r *oauthProvidersRepository) Create(ctx context.Context, oauthProviders *Oauthproviders) error {
	return r.db.WithContext(ctx).Create(oauthProviders).Error
}

// GetByID gets a oauthProviders by ID
func (r *oauthProvidersRepository) GetByID(ctx context.Context, id uint) (*Oauthproviders, error) {
	var oauthProviders Oauthproviders
	err := r.db.WithContext(ctx).First(&oauthProviders, id).Error
	if err != nil {
		return nil, err
	}
	return &oauthProviders, nil
}

// GetAll gets all oauthProviderss with pagination
func (r *oauthProvidersRepository) GetAll(ctx context.Context, limit, offset int) ([]Oauthproviders, int64, error) {
	var oauthProviderss []Oauthproviders
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Oauthproviders{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&oauthProviderss).Error
	return oauthProviderss, total, err
}

// Update updates a oauthProviders
func (r *oauthProvidersRepository) Update(ctx context.Context, oauthProviders *Oauthproviders) error {
	return r.db.WithContext(ctx).Save(oauthProviders).Error
}

// Delete deletes a oauthProviders by ID
func (r *oauthProvidersRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Oauthproviders{}, id).Error
}
