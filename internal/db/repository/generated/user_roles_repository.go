package generated

import (
	"context"
	"gorm.io/gorm"
)

// UserrolesRepository interface for user_roles operations
type UserrolesRepository interface {
	Create(ctx context.Context, userRoles *Userroles) error
	GetByID(ctx context.Context, id uint) (*Userroles, error)
	GetAll(ctx context.Context, limit, offset int) ([]Userroles, int64, error)
	Update(ctx context.Context, userRoles *Userroles) error
	Delete(ctx context.Context, id uint) error
}

// userRolesRepository implements UserrolesRepository
type userRolesRepository struct {
	db *gorm.DB
}

// NewUserrolesRepository creates a new UserrolesRepository
func NewUserrolesRepository(db *gorm.DB) UserrolesRepository {
	return &userRolesRepository{db: db}
}

// Create creates a new userRoles
func (r *userRolesRepository) Create(ctx context.Context, userRoles *Userroles) error {
	return r.db.WithContext(ctx).Create(userRoles).Error
}

// GetByID gets a userRoles by ID
func (r *userRolesRepository) GetByID(ctx context.Context, id uint) (*Userroles, error) {
	var userRoles Userroles
	err := r.db.WithContext(ctx).First(&userRoles, id).Error
	if err != nil {
		return nil, err
	}
	return &userRoles, nil
}

// GetAll gets all userRoless with pagination
func (r *userRolesRepository) GetAll(ctx context.Context, limit, offset int) ([]Userroles, int64, error) {
	var userRoless []Userroles
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Userroles{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&userRoless).Error
	return userRoless, total, err
}

// Update updates a userRoles
func (r *userRolesRepository) Update(ctx context.Context, userRoles *Userroles) error {
	return r.db.WithContext(ctx).Save(userRoles).Error
}

// Delete deletes a userRoles by ID
func (r *userRolesRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Userroles{}, id).Error
}
