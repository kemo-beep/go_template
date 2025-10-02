package generated

import (
	"context"
	"gorm.io/gorm"
)

// RolesRepository interface for roles operations
type RolesRepository interface {
	Create(ctx context.Context, roles *Roles) error
	GetByID(ctx context.Context, id uint) (*Roles, error)
	GetAll(ctx context.Context, limit, offset int) ([]Roles, int64, error)
	Update(ctx context.Context, roles *Roles) error
	Delete(ctx context.Context, id uint) error
}

// rolesRepository implements RolesRepository
type rolesRepository struct {
	db *gorm.DB
}

// NewRolesRepository creates a new RolesRepository
func NewRolesRepository(db *gorm.DB) RolesRepository {
	return &rolesRepository{db: db}
}

// Create creates a new roles
func (r *rolesRepository) Create(ctx context.Context, roles *Roles) error {
	return r.db.WithContext(ctx).Create(roles).Error
}

// GetByID gets a roles by ID
func (r *rolesRepository) GetByID(ctx context.Context, id uint) (*Roles, error) {
	var roles Roles
	err := r.db.WithContext(ctx).First(&roles, id).Error
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

// GetAll gets all roless with pagination
func (r *rolesRepository) GetAll(ctx context.Context, limit, offset int) ([]Roles, int64, error) {
	var roless []Roles
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Roles{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&roless).Error
	return roless, total, err
}

// Update updates a roles
func (r *rolesRepository) Update(ctx context.Context, roles *Roles) error {
	return r.db.WithContext(ctx).Save(roles).Error
}

// Delete deletes a roles by ID
func (r *rolesRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Roles{}, id).Error
}
