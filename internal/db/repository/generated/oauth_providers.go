package generated

import "time"


// Oauthproviders represents the oauth_providers table
type Oauthproviders struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('oauth_providers_id_seq'::regclass)" json:"id"`

	Userid uint `gorm:"not null" json:"user_id"`

	Provider string `gorm:"not null" json:"provider"`

	Provideruserid string `gorm:"not null" json:"provider_user_id"`

	Accesstoken string `gorm:"" json:"access_token"`

	Refreshtoken string `gorm:"" json:"refresh_token"`

	Tokenexpiresat time.Time `gorm:"" json:"token_expires_at"`

	Profiledata string `gorm:"" json:"profile_data"`

	Createdat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	Updatedat time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

}

// TableName returns the table name for Oauthproviders
func (Oauthproviders) TableName() string {
	return "oauth_providers"
}
