package generated

import (
	"context"
	"gorm.io/gorm"
)

// TesttableRepository interface for test_table operations
type TesttableRepository interface {
	Create(ctx context.Context, testTable *Testtable) error
	GetByID(ctx context.Context, id uint) (*Testtable, error)
	GetAll(ctx context.Context, limit, offset int) ([]Testtable, int64, error)
	Update(ctx context.Context, testTable *Testtable) error
	Delete(ctx context.Context, id uint) error
}

// testTableRepository implements TesttableRepository
type testTableRepository struct {
	db *gorm.DB
}

// NewTesttableRepository creates a new TesttableRepository
func NewTesttableRepository(db *gorm.DB) TesttableRepository {
	return &testTableRepository{db: db}
}

// Create creates a new testTable
func (r *testTableRepository) Create(ctx context.Context, testTable *Testtable) error {
	return r.db.WithContext(ctx).Create(testTable).Error
}

// GetByID gets a testTable by ID
func (r *testTableRepository) GetByID(ctx context.Context, id uint) (*Testtable, error) {
	var testTable Testtable
	err := r.db.WithContext(ctx).First(&testTable, id).Error
	if err != nil {
		return nil, err
	}
	return &testTable, nil
}

// GetAll gets all testTables with pagination
func (r *testTableRepository) GetAll(ctx context.Context, limit, offset int) ([]Testtable, int64, error) {
	var testTables []Testtable
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&Testtable{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&testTables).Error
	return testTables, total, err
}

// Update updates a testTable
func (r *testTableRepository) Update(ctx context.Context, testTable *Testtable) error {
	return r.db.WithContext(ctx).Save(testTable).Error
}

// Delete deletes a testTable by ID
func (r *testTableRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Testtable{}, id).Error
}
