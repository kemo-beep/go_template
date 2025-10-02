package generated

import "time"


// Roles represents the roles table
type Roles struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('roles_id_seq'::regclass)" json:"id"`

	Name string `gorm:"not null" json:"name"`

	Description string `gorm:"" json:"description"`

	Createdat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	Updatedat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

}

// TableName returns the table name for Roles
func (Roles) TableName() string {
	return "roles"
}
