package models

type KMSKey struct {
	ID   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
}

// TableName overrides the table name used by User to `profiles`
func (KMSKey) TableName() string {
	return "kms_keys"
}
