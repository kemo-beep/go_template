package generated

import (
	"context"
	"gorm.io/gorm"
)

// SessionsRepository interface for sessions operations
type SessionsRepository interface {
	Create(ctx context.Context, sessions *Sessions) error
	GetByID(ctx context.Context, id uint) (*Sessions, error)
	GetAll(ctx context.Context, limit, offset int) ([]Sessions, int64, error)
	Update(ctx context.Context, sessions *Sessions) error
	Delete(ctx context.Context, id uint) error
}

// sessionsRepository implements SessionsRepository
type sessionsRepository struct {
	db *gorm.DB
}

// NewSessionsRepository creates a new SessionsRepository
func NewSessionsRepository(db *gorm.DB) SessionsRepository {
	return &sessionsRepository{db: db}
}

// Create creates a new sessions
func (r *sessionsRepository) Create(ctx context.Context, sessions *Sessions) error {
	return r.db.WithContext(ctx).Create(sessions).Error
}

// GetByID gets a sessions by ID
func (r *sessionsRepository) GetByID(ctx context.Context, id uint) (*Sessions, error) {
	var sessions Sessions
	err := r.db.WithContext(ctx).First(&sessions, id).Error
	if err != nil {
		return nil, err
	}
	return &sessions, nil
}

// GetAll gets all sessionss with pagination
func (r *sessionsRepository) GetAll(ctx context.Context, limit, offset int) ([]Sessions, int64, error) {
	var sessionss []Sessions
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Sessions{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&sessionss).Error
	return sessionss, total, err
}

// Update updates a sessions
func (r *sessionsRepository) Update(ctx context.Context, sessions *Sessions) error {
	return r.db.WithContext(ctx).Save(sessions).Error
}

// Delete deletes a sessions by ID
func (r *sessionsRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Sessions{}, id).Error
}
