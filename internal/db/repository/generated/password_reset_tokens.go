package generated

import "time"


// Passwordresettokens represents the password_reset_tokens table
type Passwordresettokens struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('password_reset_tokens_id_seq'::regclass)" json:"id"`

	Userid uint `gorm:"not null" json:"user_id"`

	Token string `gorm:"not null" json:"token"`

	Expiresat time.Time `gorm:"not null" json:"expires_at"`

	Used bool `gorm:"default:false" json:"used"`

	Usedat time.Time `gorm:"" json:"used_at"`

	Createdat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

}

// TableName returns the table name for Passwordresettokens
func (Passwordresettokens) TableName() string {
	return "password_reset_tokens"
}
