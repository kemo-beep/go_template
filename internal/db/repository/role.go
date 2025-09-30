package repository

import (
	"context"

	"gorm.io/gorm"
)

// RoleRepository defines role-related database operations
type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	GetByID(ctx context.Context, id uint) (*Role, error)
	GetByName(ctx context.Context, name string) (*Role, error)
	List(ctx context.Context, limit, offset int) ([]*Role, int64, error)
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id uint) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]Permission, error)
	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]Role, error)
	AssignRoleToUser(ctx context.Context, userID, roleID uint, assignedBy *uint) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error
}

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) GetByID(ctx context.Context, id uint) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) List(ctx context.Context, limit, offset int) ([]*Role, int64, error) {
	var roles []*Role
	var total int64

	if err := r.db.WithContext(ctx).Model(&Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Permissions").
		Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *roleRepository) Update(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Role{}, id).Error
}

func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]Permission, error) {
	var permissions []Permission
	if err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *roleRepository) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	// First, remove existing permissions
	if err := r.db.WithContext(ctx).Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}

	// Then add new permissions
	for _, permID := range permissionIDs {
		if err := r.db.WithContext(ctx).Exec(
			"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
			roleID, permID,
		).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uint) ([]Role, error) {
	var roles []Role
	if err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Preload("Permissions").
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) AssignRoleToUser(ctx context.Context, userID, roleID uint, assignedBy *uint) error {
	userRole := map[string]interface{}{
		"user_id": userID,
		"role_id": roleID,
	}
	if assignedBy != nil {
		userRole["assigned_by"] = *assignedBy
	}

	return r.db.WithContext(ctx).Table("user_roles").Create(userRole).Error
}

func (r *roleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM user_roles WHERE user_id = ? AND role_id = ?",
		userID, roleID,
	).Error
}

// PermissionRepository defines permission-related database operations
type PermissionRepository interface {
	Create(ctx context.Context, permission *Permission) error
	GetByID(ctx context.Context, id uint) (*Permission, error)
	GetByName(ctx context.Context, name string) (*Permission, error)
	List(ctx context.Context, limit, offset int) ([]*Permission, int64, error)
	Update(ctx context.Context, permission *Permission) error
	Delete(ctx context.Context, id uint) error
	CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(ctx context.Context, permission *Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *permissionRepository) GetByID(ctx context.Context, id uint) (*Permission, error) {
	var permission Permission
	if err := r.db.WithContext(ctx).First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) GetByName(ctx context.Context, name string) (*Permission, error) {
	var permission Permission
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) List(ctx context.Context, limit, offset int) ([]*Permission, int64, error) {
	var permissions []*Permission
	var total int64

	if err := r.db.WithContext(ctx).Model(&Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

func (r *permissionRepository) Update(ctx context.Context, permission *Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Permission{}, id).Error
}

func (r *permissionRepository) CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("permissions.resource = ? AND permissions.action = ?", resource, action).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *permissionRepository) GetUserPermissions(ctx context.Context, userID uint) ([]Permission, error) {
	var permissions []Permission
	err := r.db.WithContext(ctx).
		Distinct().
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&permissions).Error

	return permissions, err
}
