package generated


// Testtable represents the test_table table
type Testtable struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('test_table_id_seq'::regclass)" json:"id"`

	Name string `gorm:"" json:"name"`

	Familyname string `gorm:"" json:"family_name"`

	Prefrence string `gorm:"" json:"prefrence"`

	Preferences string `gorm:"" json:"preferences"`

	Location string `gorm:"" json:"location"`

}

// TableName returns the table name for Testtable
func (Testtable) TableName() string {
	return "test_table"
}
