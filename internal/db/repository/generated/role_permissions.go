package generated

import "time"


// Rolepermissions represents the role_permissions table
type Rolepermissions struct {

	Roleid uint `gorm:"primaryKey;not null" json:"role_id"`

	Permissionid uint `gorm:"primaryKey;not null" json:"permission_id"`

	Createdat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

}

// TableName returns the table name for Rolepermissions
func (Rolepermissions) TableName() string {
	return "role_permissions"
}
