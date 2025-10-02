package generated

import (
	"context"
	"gorm.io/gorm"
)

// RolepermissionsRepository interface for role_permissions operations
type RolepermissionsRepository interface {
	Create(ctx context.Context, rolePermissions *Rolepermissions) error
	GetByID(ctx context.Context, id uint) (*Rolepermissions, error)
	GetAll(ctx context.Context, limit, offset int) ([]Rolepermissions, int64, error)
	Update(ctx context.Context, rolePermissions *Rolepermissions) error
	Delete(ctx context.Context, id uint) error
}

// rolePermissionsRepository implements RolepermissionsRepository
type rolePermissionsRepository struct {
	db *gorm.DB
}

// NewRolepermissionsRepository creates a new RolepermissionsRepository
func NewRolepermissionsRepository(db *gorm.DB) RolepermissionsRepository {
	return &rolePermissionsRepository{db: db}
}

// Create creates a new rolePermissions
func (r *rolePermissionsRepository) Create(ctx context.Context, rolePermissions *Rolepermissions) error {
	return r.db.WithContext(ctx).Create(rolePermissions).Error
}

// GetByID gets a rolePermissions by ID
func (r *rolePermissionsRepository) GetByID(ctx context.Context, id uint) (*Rolepermissions, error) {
	var rolePermissions Rolepermissions
	err := r.db.WithContext(ctx).First(&rolePermissions, id).Error
	if err != nil {
		return nil, err
	}
	return &rolePermissions, nil
}

// GetAll gets all rolePermissionss with pagination
func (r *rolePermissionsRepository) GetAll(ctx context.Context, limit, offset int) ([]Rolepermissions, int64, error) {
	var rolePermissionss []Rolepermissions
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Rolepermissions{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&rolePermissionss).Error
	return rolePermissionss, total, err
}

// Update updates a rolePermissions
func (r *rolePermissionsRepository) Update(ctx context.Context, rolePermissions *Rolepermissions) error {
	return r.db.WithContext(ctx).Save(rolePermissions).Error
}

// Delete deletes a rolePermissions by ID
func (r *rolePermissionsRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Rolepermissions{}, id).Error
}
