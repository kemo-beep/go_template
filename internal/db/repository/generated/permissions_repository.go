package generated

import (
	"context"
	"gorm.io/gorm"
)

// PermissionsRepository interface for permissions operations
type PermissionsRepository interface {
	Create(ctx context.Context, permissions *Permissions) error
	GetByID(ctx context.Context, id uint) (*Permissions, error)
	GetAll(ctx context.Context, limit, offset int) ([]Permissions, int64, error)
	Update(ctx context.Context, permissions *Permissions) error
	Delete(ctx context.Context, id uint) error
}

// permissionsRepository implements PermissionsRepository
type permissionsRepository struct {
	db *gorm.DB
}

// NewPermissionsRepository creates a new PermissionsRepository
func NewPermissionsRepository(db *gorm.DB) PermissionsRepository {
	return &permissionsRepository{db: db}
}

// Create creates a new permissions
func (r *permissionsRepository) Create(ctx context.Context, permissions *Permissions) error {
	return r.db.WithContext(ctx).Create(permissions).Error
}

// GetByID gets a permissions by ID
func (r *permissionsRepository) GetByID(ctx context.Context, id uint) (*Permissions, error) {
	var permissions Permissions
	err := r.db.WithContext(ctx).First(&permissions, id).Error
	if err != nil {
		return nil, err
	}
	return &permissions, nil
}

// GetAll gets all permissionss with pagination
func (r *permissionsRepository) GetAll(ctx context.Context, limit, offset int) ([]Permissions, int64, error) {
	var permissionss []Permissions
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Permissions{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&permissionss).Error
	return permissionss, total, err
}

// Update updates a permissions
func (r *permissionsRepository) Update(ctx context.Context, permissions *Permissions) error {
	return r.db.WithContext(ctx).Save(permissions).Error
}

// Delete deletes a permissions by ID
func (r *permissionsRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Permissions{}, id).Error
}
