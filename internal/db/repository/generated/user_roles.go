package generated

import "time"


// Userroles represents the user_roles table
type Userroles struct {

	Userid uint `gorm:"primaryKey;not null" json:"user_id"`

	Roleid uint `gorm:"primaryKey;not null" json:"role_id"`

	Assignedat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`

	Assignedby uint `gorm:"" json:"assigned_by"`

}

// TableName returns the table name for Userroles
func (Userroles) TableName() string {
	return "user_roles"
}
