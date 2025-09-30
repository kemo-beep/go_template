package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// fileRepository implements FileRepository interface
type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository creates a new file repository
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

// Create creates a new file record
func (r *fileRepository) Create(ctx context.Context, file *File) error {
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	return nil
}

// GetByID retrieves a file by ID
func (r *fileRepository) GetByID(ctx context.Context, id uint) (*File, error) {
	var file File
	if err := r.db.WithContext(ctx).Preload("User").First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return &file, nil
}

// GetByR2Key retrieves a file by R2 key
func (r *fileRepository) GetByR2Key(ctx context.Context, r2Key string) (*File, error) {
	var file File
	if err := r.db.WithContext(ctx).Where("r2_key = ?", r2Key).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to get file by R2 key: %w", err)
	}
	return &file, nil
}

// GetByUserID retrieves files by user ID with pagination
func (r *fileRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*File, error) {
	var files []*File
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&files).Error; err != nil {
		return nil, fmt.Errorf("failed to get files by user ID: %w", err)
	}
	return files, nil
}

// Update updates an existing file
func (r *fileRepository) Update(ctx context.Context, file *File) error {
	if err := r.db.WithContext(ctx).Save(file).Error; err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}
	return nil
}

// Delete soft deletes a file
func (r *fileRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&File{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// DeleteByR2Key deletes a file by R2 key
func (r *fileRepository) DeleteByR2Key(ctx context.Context, r2Key string) error {
	if err := r.db.WithContext(ctx).Where("r2_key = ?", r2Key).Delete(&File{}).Error; err != nil {
		return fmt.Errorf("failed to delete file by R2 key: %w", err)
	}
	return nil
}
