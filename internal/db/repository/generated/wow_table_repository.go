package generated

import (
	"context"
	"gorm.io/gorm"
)

// WowtableRepository interface for wow_table operations
type WowtableRepository interface {
	Create(ctx context.Context, wowTable *Wowtable) error
	GetByID(ctx context.Context, id uint) (*Wowtable, error)
	GetAll(ctx context.Context, limit, offset int) ([]Wowtable, int64, error)
	Update(ctx context.Context, wowTable *Wowtable) error
	Delete(ctx context.Context, id uint) error
}

// wowTableRepository implements WowtableRepository
type wowTableRepository struct {
	db *gorm.DB
}

// NewWowtableRepository creates a new WowtableRepository
func NewWowtableRepository(db *gorm.DB) WowtableRepository {
	return &wowTableRepository{db: db}
}

// Create creates a new wowTable
func (r *wowTableRepository) Create(ctx context.Context, wowTable *Wowtable) error {
	return r.db.WithContext(ctx).Create(wowTable).Error
}

// GetByID gets a wowTable by ID
func (r *wowTableRepository) GetByID(ctx context.Context, id uint) (*Wowtable, error) {
	var wowTable Wowtable
	err := r.db.WithContext(ctx).First(&wowTable, id).Error
	if err != nil {
		return nil, err
	}
	return &wowTable, nil
}

// GetAll gets all wowTables with pagination
func (r *wowTableRepository) GetAll(ctx context.Context, limit, offset int) ([]Wowtable, int64, error) {
	var wowTables []Wowtable
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Wowtable{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&wowTables).Error
	return wowTables, total, err
}

// Update updates a wowTable
func (r *wowTableRepository) Update(ctx context.Context, wowTable *Wowtable) error {
	return r.db.WithContext(ctx).Save(wowTable).Error
}

// Delete deletes a wowTable by ID
func (r *wowTableRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Wowtable{}, id).Error
}
