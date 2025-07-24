package models

import "time"

type KMSKey struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Alias      string         `gorm:"type:varchar(255);not null" json:"name"`
	Algorithm  string         `json:"algorithm"`
	Size       int            `json:"size"`
	PublicKey  string         `json:"public_key"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreationTS time.Time      `json:"creation_ts"`
}

// TableName overrides the table name used by User to `profiles`
func (KMSKey) TableName() string {
	return "kms_keys"
}
