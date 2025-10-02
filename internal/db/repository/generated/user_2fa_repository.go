package generated

import (
	"context"
	"gorm.io/gorm"
)

// User2faRepository interface for user_2fa operations
type User2faRepository interface {
	Create(ctx context.Context, user2fa *User2fa) error
	GetByID(ctx context.Context, id uint) (*User2fa, error)
	GetAll(ctx context.Context, limit, offset int) ([]User2fa, int64, error)
	Update(ctx context.Context, user2fa *User2fa) error
	Delete(ctx context.Context, id uint) error
}

// user2faRepository implements User2faRepository
type user2faRepository struct {
	db *gorm.DB
}

// NewUser2faRepository creates a new User2faRepository
func NewUser2faRepository(db *gorm.DB) User2faRepository {
	return &user2faRepository{db: db}
}

// Create creates a new user2fa
func (r *user2faRepository) Create(ctx context.Context, user2fa *User2fa) error {
	return r.db.WithContext(ctx).Create(user2fa).Error
}

// GetByID gets a user2fa by ID
func (r *user2faRepository) GetByID(ctx context.Context, id uint) (*User2fa, error) {
	var user2fa User2fa
	err := r.db.WithContext(ctx).First(&user2fa, id).Error
	if err != nil {
		return nil, err
	}
	return &user2fa, nil
}

// GetAll gets all user2fas with pagination
func (r *user2faRepository) GetAll(ctx context.Context, limit, offset int) ([]User2fa, int64, error) {
	var user2fas []User2fa
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&User2fa{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&user2fas).Error
	return user2fas, total, err
}

// Update updates a user2fa
func (r *user2faRepository) Update(ctx context.Context, user2fa *User2fa) error {
	return r.db.WithContext(ctx).Save(user2fa).Error
}

// Delete deletes a user2fa by ID
func (r *user2faRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&User2fa{}, id).Error
}
