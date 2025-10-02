package generated

import "time"


// Refreshtokens represents the refresh_tokens table
type Refreshtokens struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('refresh_tokens_id_seq'::regclass)" json:"id"`

	Userid uint `gorm:"not null" json:"user_id"`

	Token string `gorm:"not null" json:"token"`

	Expiresat time.Time `gorm:"not null" json:"expires_at"`

	Isrevoked bool `gorm:"default:false" json:"is_revoked"`

	Createdat time.Time `gorm:"" json:"created_at"`

	Updatedat time.Time `gorm:"" json:"updated_at"`

	Deletedat time.Time `gorm:"" json:"deleted_at"`

}

// TableName returns the table name for Refreshtokens
func (Refreshtokens) TableName() string {
	return "refresh_tokens"
}
