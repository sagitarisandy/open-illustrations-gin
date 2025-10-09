package services

import (
	"time"

	"open-illustrations-go/config"
	"open-illustrations-go/models"
)

func CreateCategory(name string) (*models.Category, error) {
	cat := models.Category{Name: name, Slug: slugify(name)}
	if err := config.DB.Create(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func GetCategories() ([]models.Category, error) {
	var list []models.Category
	res := config.DB.Where("deleted_at IS NULL").Find(&list)
	return list, res.Error
}

func GetCategory(id string) (*models.Category, error) {
	var c models.Category
	res := config.DB.First(&c, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &c, nil
}

func UpdateCategory(id string, name string) (*models.Category, error) {
	c, err := GetCategory(id)
	if err != nil {
		return nil, err
	}
	c.Name = name
	c.Slug = slugify(name)
	if err := config.DB.Save(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

// Legacy hard delete (kept for compatibility)
func DeleteCategory(id string) error {
	return config.DB.Delete(&models.Category{}, id).Error
}

func SoftDeleteCategory(id string) (*models.Category, error) {
	var c models.Category
	if err := config.DB.First(&c, id).Error; err != nil {
		return nil, err
	}
	if c.DeletedAt.Valid { // already deleted
		return &c, nil
	}
	if err := config.DB.Model(&c).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
