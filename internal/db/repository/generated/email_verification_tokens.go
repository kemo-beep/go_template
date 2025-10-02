package generated

import "time"


// Emailverificationtokens represents the email_verification_tokens table
type Emailverificationtokens struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('email_verification_tokens_id_seq'::regclass)" json:"id"`

	Userid uint `gorm:"not null" json:"user_id"`

	Email string `gorm:"not null" json:"email"`

	Token string `gorm:"not null" json:"token"`

	Expiresat time.Time `gorm:"not null" json:"expires_at"`

	Used bool `gorm:"default:false" json:"used"`

	Usedat time.Time `gorm:"" json:"used_at"`

	Createdat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

}

// TableName returns the table name for Emailverificationtokens
func (Emailverificationtokens) TableName() string {
	return "email_verification_tokens"
}
