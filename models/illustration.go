package models

import (
	"time"

	"gorm.io/gorm"
)

type Illustration struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	StyleID    *uint          `gorm:"column:style_id;index" json:"style_id"`
	CategoryID *uint          `gorm:"column:category_id;index" json:"category_id"`
	PackID     *uint          `gorm:"column:pack_id;index" json:"pack_id"`
	FileName   string         `gorm:"size:191;not null" json:"file_name"`
	StorageKey string         `gorm:"size:191;not null;uniqueIndex" json:"storage_key"`
	IsPremium  bool           `gorm:"index" json:"is_premium"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	CategoryRef *Category `gorm:"foreignKey:CategoryID" json:"category_ref,omitempty"`
	PackRef     *Pack     `gorm:"foreignKey:PackID" json:"pack_ref,omitempty"`
	StyleRef    *Style    `gorm:"foreignKey:StyleID" json:"style_ref,omitempty"`
}

type Category struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug          string         `gorm:"size:120;uniqueIndex" json:"slug,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Illustrations []Illustration `gorm:"foreignKey:CategoryID" json:"illustrations,omitempty"`
}

type Pack struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug          string         `gorm:"size:120;uniqueIndex" json:"slug,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Illustrations []Illustration `gorm:"foreignKey:PackID" json:"illustrations,omitempty"`
}

type Style struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug         string         `gorm:"size:120;uniqueIndex" json:"slug,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Illustration []Illustration `gorm:"foreignKey:StyleID" json:"illustrations,omitempty"`
}

// TableName explicitly keeps the pluralized form for Pack if needed (GORM would default to "packs" already, provided for clarity).
func (Pack) TableName() string { return "packs" }
