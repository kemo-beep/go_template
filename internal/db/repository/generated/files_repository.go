package generated

import (
	"context"
	"gorm.io/gorm"
)

// FilesRepository interface for files operations
type FilesRepository interface {
	Create(ctx context.Context, files *Files) error
	GetByID(ctx context.Context, id uint) (*Files, error)
	GetAll(ctx context.Context, limit, offset int) ([]Files, int64, error)
	Update(ctx context.Context, files *Files) error
	Delete(ctx context.Context, id uint) error
}

// filesRepository implements FilesRepository
type filesRepository struct {
	db *gorm.DB
}

// NewFilesRepository creates a new FilesRepository
func NewFilesRepository(db *gorm.DB) FilesRepository {
	return &filesRepository{db: db}
}

// Create creates a new files
func (r *filesRepository) Create(ctx context.Context, files *Files) error {
	return r.db.WithContext(ctx).Create(files).Error
}

// GetByID gets a files by ID
func (r *filesRepository) GetByID(ctx context.Context, id uint) (*Files, error) {
	var files Files
	err := r.db.WithContext(ctx).First(&files, id).Error
	if err != nil {
		return nil, err
	}
	return &files, nil
}

// GetAll gets all filess with pagination
func (r *filesRepository) GetAll(ctx context.Context, limit, offset int) ([]Files, int64, error) {
	var filess []Files
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Files{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&filess).Error
	return filess, total, err
}

// Update updates a files
func (r *filesRepository) Update(ctx context.Context, files *Files) error {
	return r.db.WithContext(ctx).Save(files).Error
}

// Delete deletes a files by ID
func (r *filesRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Files{}, id).Error
}
