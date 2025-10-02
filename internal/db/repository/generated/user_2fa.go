package generated

import "time"


// User2fa represents the user_2fa table
type User2fa struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('user_2fa_id_seq'::regclass)" json:"id"`

	Userid uint `gorm:"not null" json:"user_id"`

	Secret string `gorm:"not null" json:"secret"`

	Backupcodes string `gorm:"" json:"backup_codes"`

	Isenabled bool `gorm:"default:false" json:"is_enabled"`

	Enabledat time.Time `gorm:"" json:"enabled_at"`

	Lastusedat time.Time `gorm:"" json:"last_used_at"`

	Createdat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	Updatedat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

}

// TableName returns the table name for User2fa
func (User2fa) TableName() string {
	return "user_2fa"
}
