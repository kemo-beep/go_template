package generated


// Dancingtable represents the dancing_table table
type Dancingtable struct {

	Id uint `gorm:"primaryKey;not null;default:nextval('dancing_table_id_seq'::regclass)" json:"id"`

	Frequency string `gorm:"" json:"frequency"`

}

// TableName returns the table name for Dancingtable
func (Dancingtable) TableName() string {
	return "dancing_table"
}
