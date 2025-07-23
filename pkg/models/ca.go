package models

type CACertificate struct {
	ID   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
}

// TableName overrides the table name used by User to `profiles`
func (CACertificate) TableName() string {
	return "cas"
}
