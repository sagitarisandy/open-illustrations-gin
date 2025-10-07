package models

import (
	"time"

	"gorm.io/gorm"
)

// Illustration represents a single illustration asset stored in MinIO and indexed in MySQL.
// We keep a legacy Category string for simple filtering while also allowing relational Category/Pack.
type Illustration struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	CategoryID *uint          `gorm:"column:category_id;index" json:"category_id,omitempty"`
	PackID     *uint          `gorm:"column:pack_id;index" json:"pack_id,omitempty"`
	FileName   string         `gorm:"size:191;not null" json:"file_name"`
	StorageKey string         `gorm:"size:191;not null;uniqueIndex" json:"storage_key"`
	IsPremium  bool           `gorm:"index" json:"is_premium"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Associations (omit from JSON by default to avoid heavy nesting unless explicitly preloaded)
	CategoryRef *Category `gorm:"foreignKey:CategoryID" json:"category_ref,omitempty"`
	PackRef     *Pack     `gorm:"foreignKey:PackID" json:"pack_ref,omitempty"`
}

// Category groups illustrations; unique by Name.
type Category struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug          string         `gorm:"size:120;uniqueIndex" json:"slug,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Illustrations []Illustration `gorm:"foreignKey:CategoryID" json:"illustrations,omitempty"`
}

// Pack (formerly Packs) groups a collection of illustrations into a themed bundle.
type Pack struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug          string         `gorm:"size:120;uniqueIndex" json:"slug,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Illustrations []Illustration `gorm:"foreignKey:PackID" json:"illustrations,omitempty"`
}

// TableName explicitly keeps the pluralized form for Pack if needed (GORM would default to "packs" already, provided for clarity).
func (Pack) TableName() string { return "packs" }
