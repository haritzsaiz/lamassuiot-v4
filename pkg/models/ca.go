package models

type CACertificate struct {
	ID string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
}
