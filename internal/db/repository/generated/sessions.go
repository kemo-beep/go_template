package generated

import "time"


// Sessions represents the sessions table
type Sessions struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('sessions_id_seq'::regclass)" json:"id"`

	Userid uint `gorm:"not null" json:"user_id"`

	Token string `gorm:"not null" json:"token"`

	Refreshtoken string `gorm:"" json:"refresh_token"`

	Deviceinfo string `gorm:"" json:"device_info"`

	Ipaddress string `gorm:"" json:"ip_address"`

	Useragent string `gorm:"" json:"user_agent"`

	Isactive bool `gorm:"default:true" json:"is_active"`

	Expiresat time.Time `gorm:"not null" json:"expires_at"`

	Createdat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	Lastusedat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_used_at"`

}

// TableName returns the table name for Sessions
func (Sessions) TableName() string {
	return "sessions"
}
