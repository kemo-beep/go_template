package repository

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// JSONB is a custom type for PostgreSQL JSONB
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = result
	return nil
}

// User represents a user in the system
type User struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	Email               string         `json:"email" gorm:"uniqueIndex;not null"`
	Password            string         `json:"-" gorm:"not null"`
	Name                string         `json:"name" gorm:"not null"`
	IsActive            bool           `json:"is_active" gorm:"default:true"`
	IsAdmin             bool           `json:"is_admin" gorm:"default:false"`
	EmailVerified       bool           `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt     *time.Time     `json:"email_verified_at,omitempty"`
	LastLoginAt         *time.Time     `json:"last_login_at,omitempty"`
	FailedLoginAttempts int            `json:"-" gorm:"default:0"`
	LockedUntil         *time.Time     `json:"locked_until,omitempty"`
	Metadata            JSONB          `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

// File represents a file stored in R2
type File struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	FileName  string         `json:"file_name" gorm:"not null"`
	FileSize  int64          `json:"file_size" gorm:"not null"`
	FileType  string         `json:"file_type" gorm:"not null"`
	R2Key     string         `json:"r2_key" gorm:"uniqueIndex;not null"`
	R2URL     string         `json:"r2_url" gorm:"not null"`
	IsPublic  bool           `json:"is_public" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
}

// RefreshToken represents a JWT refresh token
type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	IsRevoked bool           `json:"is_revoked" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
}

// Role represents a user role in the system
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission represents a system permission
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
	Resource    string    `gorm:"not null" json:"resource"`
	Action      string    `gorm:"not null" json:"action"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// TableName returns the table name for File
func (File) TableName() string {
	return "files"
}

// TableName returns the table name for RefreshToken
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// TableName returns the table name for Role
func (Role) TableName() string {
	return "roles"
}

// TableName returns the table name for Permission
func (Permission) TableName() string {
	return "permissions"
}
