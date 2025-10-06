package models

import "time"

type Illustration struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"size:200;not null" json:"title"`
	Category   string    `gorm:"size:100;not null;index" json:"category"`
	FileName   string    `gorm:"size:191;not null" json:"file_name"`
	StorageKey string    `gorm:"size:191;not null;uniqueIndex" json:"storage_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
