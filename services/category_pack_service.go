package services

import (
	"fmt"
	"strings"
	"time"

	"open-illustrations-go/config"
	"open-illustrations-go/models"
)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

// ---- Category ----
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

func DeleteCategory(id string) error { // legacy hard delete (unused now)
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

// ---- Pack ----
func CreatePack(name string) (*models.Pack, error) {
	p := models.Pack{Name: name, Slug: slugify(name)}
	if err := config.DB.Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func GetPacks() ([]models.Pack, error) {
	var list []models.Pack
	res := config.DB.Where("deleted_at IS NULL").Find(&list)
	return list, res.Error
}

func GetPack(id string) (*models.Pack, error) {
	var p models.Pack
	res := config.DB.First(&p, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &p, nil
}

func UpdatePack(id string, name string) (*models.Pack, error) {
	p, err := GetPack(id)
	if err != nil {
		return nil, err
	}
	p.Name = name
	p.Slug = slugify(name)
	if err := config.DB.Save(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func DeletePack(id string) error { // legacy hard delete
	return config.DB.Delete(&models.Pack{}, id).Error
}

func SoftDeletePack(id string) (*models.Pack, error) {
	var p models.Pack
	if err := config.DB.First(&p, id).Error; err != nil {
		return nil, err
	}
	if p.DeletedAt.Valid {
		return &p, nil
	}
	if err := config.DB.Model(&p).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// Generate a simple pseudo archive name for pack download (caller will stream objects)
func PackArchiveFileName(p *models.Pack) string {
	return fmt.Sprintf("pack-%s-%d-%s.zip", p.Slug, p.ID, time.Now().Format("20060102"))
}
