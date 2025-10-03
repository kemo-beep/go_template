package generated

import (
	"context"
	"gorm.io/gorm"
)

// DancingtableRepository interface for dancing_table operations
type DancingtableRepository interface {
	Create(ctx context.Context, dancingTable *Dancingtable) error
	GetByID(ctx context.Context, id uint) (*Dancingtable, error)
	GetAll(ctx context.Context, limit, offset int) ([]Dancingtable, int64, error)
	Update(ctx context.Context, dancingTable *Dancingtable) error
	Delete(ctx context.Context, id uint) error
}

// dancingTableRepository implements DancingtableRepository
type dancingTableRepository struct {
	db *gorm.DB
}

// NewDancingtableRepository creates a new DancingtableRepository
func NewDancingtableRepository(db *gorm.DB) DancingtableRepository {
	return &dancingTableRepository{db: db}
}

// Create creates a new dancingTable
func (r *dancingTableRepository) Create(ctx context.Context, dancingTable *Dancingtable) error {
	return r.db.WithContext(ctx).Create(dancingTable).Error
}

// GetByID gets a dancingTable by ID
func (r *dancingTableRepository) GetByID(ctx context.Context, id uint) (*Dancingtable, error) {
	var dancingTable Dancingtable
	err := r.db.WithContext(ctx).First(&dancingTable, id).Error
	if err != nil {
		return nil, err
	}
	return &dancingTable, nil
}

// GetAll gets all dancingTables with pagination
func (r *dancingTableRepository) GetAll(ctx context.Context, limit, offset int) ([]Dancingtable, int64, error) {
	var dancingTables []Dancingtable
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Dancingtable{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&dancingTables).Error
	return dancingTables, total, err
}

// Update updates a dancingTable
func (r *dancingTableRepository) Update(ctx context.Context, dancingTable *Dancingtable) error {
	return r.db.WithContext(ctx).Save(dancingTable).Error
}

// Delete deletes a dancingTable by ID
func (r *dancingTableRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Dancingtable{}, id).Error
}
