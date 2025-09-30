package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

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

// UserRole represents the user-role assignment
type UserRole struct {
	UserID     uint      `gorm:"primaryKey" json:"user_id"`
	RoleID     uint      `gorm:"primaryKey" json:"role_id"`
	AssignedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`
	AssignedBy *uint     `json:"assigned_by,omitempty"`
	Role       Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

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

// Session represents a user session
type Session struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	Token        string    `gorm:"uniqueIndex;not null" json:"-"`
	RefreshToken string    `gorm:"uniqueIndex" json:"-"`
	DeviceInfo   JSONB     `gorm:"type:jsonb" json:"device_info,omitempty"`
	IPAddress    string    `json:"ip_address,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	IsActive     bool      `gorm:"default:true;index" json:"is_active"`
	ExpiresAt    time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	LastUsedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_used_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// OAuthProvider represents OAuth provider linkage
type OAuthProvider struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"not null;index" json:"user_id"`
	Provider       string     `gorm:"not null;index" json:"provider"` // google, github, etc.
	ProviderUserID string     `gorm:"not null" json:"provider_user_id"`
	AccessToken    string     `json:"-"`
	RefreshToken   string     `json:"-"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	ProfileData    JSONB      `gorm:"type:jsonb" json:"profile_data,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// User2FA represents two-factor authentication data
type User2FA struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	Secret      string     `gorm:"not null" json:"-"`
	BackupCodes JSONB      `gorm:"type:jsonb" json:"-"`
	IsEnabled   bool       `gorm:"default:false" json:"is_enabled"`
	EnabledAt   *time.Time `json:"enabled_at,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// PasswordResetToken represents password reset token
type PasswordResetToken struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	Token     string     `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	Used      bool       `gorm:"default:false" json:"used"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// EmailVerificationToken represents email verification token
type EmailVerificationToken struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	Email     string     `gorm:"not null" json:"email"`
	Token     string     `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	Used      bool       `gorm:"default:false" json:"used"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// APIKey represents an API key for service accounts
type APIKey struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     *uint      `gorm:"index" json:"user_id,omitempty"`
	Name       string     `gorm:"not null" json:"name"`
	KeyHash    string     `gorm:"uniqueIndex;not null" json:"-"`
	Prefix     string     `gorm:"not null" json:"prefix"`
	Scopes     JSONB      `gorm:"type:jsonb" json:"scopes,omitempty"`
	RateLimit  int        `gorm:"default:1000" json:"rate_limit"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     *uint     `gorm:"index" json:"user_id,omitempty"`
	Action     string    `gorm:"not null;index" json:"action"`
	Resource   string    `json:"resource,omitempty"`
	ResourceID string    `json:"resource_id,omitempty"`
	IPAddress  string    `json:"ip_address,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
	Metadata   JSONB     `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName overrides for GORM
func (Role) TableName() string                   { return "roles" }
func (Permission) TableName() string             { return "permissions" }
func (UserRole) TableName() string               { return "user_roles" }
func (Session) TableName() string                { return "sessions" }
func (OAuthProvider) TableName() string          { return "oauth_providers" }
func (User2FA) TableName() string                { return "user_2fa" }
func (PasswordResetToken) TableName() string     { return "password_reset_tokens" }
func (EmailVerificationToken) TableName() string { return "email_verification_tokens" }
func (APIKey) TableName() string                 { return "api_keys" }
func (AuditLog) TableName() string               { return "audit_logs" }
