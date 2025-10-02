package generated

import "time"


// Permissions represents the permissions table
type Permissions struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('permissions_id_seq'::regclass)" json:"id"`

	Name string `gorm:"not null" json:"name"`

	Description string `gorm:"" json:"description"`

	Resource string `gorm:"not null" json:"resource"`

	Action string `gorm:"not null" json:"action"`

	Createdat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	Updatedat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

}

// TableName returns the table name for Permissions
func (Permissions) TableName() string {
	return "permissions"
}
