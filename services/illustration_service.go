package services

import (
	"context"
	"open-illustrations-go/config"
	"open-illustrations-go/models"
	"time"
)

// var illustrations = []models.Illustration {
// 	{ID: "1", Title: "Team Work", Category: "Business", FileName: "teamwork.svg"},
// 	{ID: "2", Title: "Coding", Category: "Technology", FileName: "coding.svg"},
// }

func GetIllustrations() ([]models.Illustration, error) {
	var illustrations []models.Illustration
	result := config.DB.Find(&illustrations)
	return illustrations, result.Error
}

func GetIllustration(id string) (*models.Illustration, error) {
	var illustration models.Illustration
	result := config.DB.First(&illustration, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &illustration, nil
}

func CreateIllustration(ill *models.Illustration) error {
	return config.DB.Create(ill).Error
}

func DeleteIllustration(id string) error {
	return config.DB.Delete(&models.Illustration{}, id).Error
}

// Generate presigned download URL from MinIO
func GetDownloadURL(fileName string, duration time.Duration) (string, error) {
	ctx := context.Background()
	reqParams := make(map[string]string)

	url, err := config.MinioClient.PresignedGetObject(ctx, config.BucketName, fileName, duration, reqParams)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
