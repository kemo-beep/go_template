package generated


// Wowtable represents the wow_table table
type Wowtable struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('wow_table_id_seq'::regclass)" json:"id"`

	Swim string `gorm:"" json:"swim"`

}

// TableName returns the table name for Wowtable
func (Wowtable) TableName() string {
	return "wow_table"
}
