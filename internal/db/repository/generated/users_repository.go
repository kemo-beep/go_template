package generated

import (
	"context"
	"gorm.io/gorm"
)

// UsersRepository interface for users operations
type UsersRepository interface {
	Create(ctx context.Context, users *Users) error
	GetByID(ctx context.Context, id uint) (*Users, error)
	GetAll(ctx context.Context, limit, offset int) ([]Users, int64, error)
	Update(ctx context.Context, users *Users) error
	Delete(ctx context.Context, id uint) error
}

// usersRepository implements UsersRepository
type usersRepository struct {
	db *gorm.DB
}

// NewUsersRepository creates a new UsersRepository
func NewUsersRepository(db *gorm.DB) UsersRepository {
	return &usersRepository{db: db}
}

// Create creates a new users
func (r *usersRepository) Create(ctx context.Context, users *Users) error {
	return r.db.WithContext(ctx).Create(users).Error
}

// GetByID gets a users by ID
func (r *usersRepository) GetByID(ctx context.Context, id uint) (*Users, error) {
	var users Users
	err := r.db.WithContext(ctx).First(&users, id).Error
	if err != nil {
		return nil, err
	}
	return &users, nil
}

// GetAll gets all userss with pagination
func (r *usersRepository) GetAll(ctx context.Context, limit, offset int) ([]Users, int64, error) {
	var userss []Users
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Users{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&userss).Error
	return userss, total, err
}

// Update updates a users
func (r *usersRepository) Update(ctx context.Context, users *Users) error {
	return r.db.WithContext(ctx).Save(users).Error
}

// Delete deletes a users by ID
func (r *usersRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Users{}, id).Error
}
