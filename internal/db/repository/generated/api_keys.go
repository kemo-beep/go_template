package generated

import "time"


// Apikeys represents the api_keys table
type Apikeys struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('api_keys_id_seq'::regclass)" json:"id"`

	Userid uint `gorm:"" json:"user_id"`

	Name string `gorm:"not null" json:"name"`

	Keyhash string `gorm:"not null" json:"key_hash"`

	Prefix string `gorm:"not null" json:"prefix"`

	Scopes string `gorm:"" json:"scopes"`

	Ratelimit uint `gorm:"default:1000" json:"rate_limit"`

	Isactive bool `gorm:"default:true" json:"is_active"`

	Lastusedat time.Time `gorm:"" json:"last_used_at"`

	Expiresat time.Time `gorm:"" json:"expires_at"`

	Createdat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	Updatedat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

}

// TableName returns the table name for Apikeys
func (Apikeys) TableName() string {
	return "api_keys"
}
