package generated

import "time"


// Files represents the files table
type Files struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('files_id_seq'::regclass)" json:"id"`

	Userid uint `gorm:"not null" json:"user_id"`

	Filename string `gorm:"not null" json:"file_name"`

	Filesize uint `gorm:"not null" json:"file_size"`

	Filetype string `gorm:"not null" json:"file_type"`

	R2key string `gorm:"not null" json:"r2_key"`

	R2url string `gorm:"not null" json:"r2_url"`

	Ispublic bool `gorm:"default:false" json:"is_public"`

	Createdat time.Time `gorm:"" json:"created_at"`

	Updatedat time.Time `gorm:"" json:"updated_at"`

	Deletedat time.Time `gorm:"" json:"deleted_at"`

}

// TableName returns the table name for Files
func (Files) TableName() string {
	return "files"
}
