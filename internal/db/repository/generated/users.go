package generated

import "time"


// Users represents the users table
type Users struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('users_id_seq'::regclass)" json:"id"`

	Email string `gorm:"not null" json:"email"`

	Password string `gorm:"not null" json:"password"`

	Name string `gorm:"not null" json:"name"`

	Isactive bool `gorm:"default:true" json:"is_active"`

	Isadmin bool `gorm:"default:false" json:"is_admin"`

	Createdat time.Time `gorm:"" json:"created_at"`

	Updatedat time.Time `gorm:"" json:"updated_at"`

	Deletedat time.Time `gorm:"" json:"deleted_at"`

	Emailverified bool `gorm:"default:false" json:"email_verified"`

	Emailverifiedat time.Time `gorm:"" json:"email_verified_at"`

	Lastloginat time.Time `gorm:"" json:"last_login_at"`

	Failedloginattempts uint `gorm:"default:0" json:"failed_login_attempts"`

	Lockeduntil time.Time `gorm:"" json:"locked_until"`

	Metadata string `gorm:"" json:"metadata"`

	Nickname string `gorm:"" json:"nickname"`

	Bio string `gorm:"" json:"bio"`

}

// TableName returns the table name for Users
func (Users) TableName() string {
	return "users"
}
