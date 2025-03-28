package models

import "gorm.io/gorm"

// CreateLocationInfoGINIndex creates the GIN index on
// the 'info' column explicitly specifying jsonb
func CreateLocationInfoGINIndex(db *gorm.DB) error {
	return db.Exec(`CREATE INDEX IF NOT EXISTS idx_locations_info 
	ON locations USING gin (info jsonb_ops);`).Error
}
