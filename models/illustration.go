package models

import "gorm.io/gorm"

type Illustration struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title    string `json:"title" gorm:"type:varchar(191);not null"`
	Category string `json:"category" gorm:"type:varchar(191);not null"`
	FileName string `json:"file_name" gorm:"type:varchar(191);not null;index"`
}
