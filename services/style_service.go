package services

import (
	"strings"
	"time"

	"open-illustrations-go/config"
	"open-illustrations-go/models"
)

func slugifyStyle(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func CreateStyle(name string) (*models.Style, error) {
	s := models.Style{Name: name, Slug: slugifyStyle(name)}
	if err := config.DB.Create(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func GetStyles() ([]models.Style, error) {
	var list []models.Style
	res := config.DB.Where("deleted_at IS NULL").Find(&list)
	return list, res.Error
}

func GetStyle(id string) (*models.Style, error) {
	var s models.Style
	res := config.DB.First(&s, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &s, nil
}

func UpdateStyle(id string, name string) (*models.Style, error) {
	s, err := GetStyle(id)
	if err != nil {
		return nil, err
	}
	s.Name = name
	s.Slug = slugifyStyle(name)
	if err := config.DB.Save(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func SoftDeleteStyle(id string) (*models.Style, error) {
	var s models.Style
	if err := config.DB.First(&s, id).Error; err != nil {
		return nil, err
	}
	if s.DeletedAt.Valid {
		return &s, nil
	}
	if err := config.DB.Model(&s).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, err
	}
	return &s, nil
}
