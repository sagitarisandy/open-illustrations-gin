package services

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"open-illustrations-go/config"
	"open-illustrations-go/models"

	"github.com/minio/minio-go/v7"
)

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

func minioObjectExists(objectName string) (bool, error) {
	_, err := config.MinioClient.StatObject(context.Background(), config.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func CreateIllustration(ill *models.Illustration) error {
	ok, err := minioObjectExists(ill.FileName)
	if err != nil {
		return fmt.Errorf("minio check failed: %w", err)
	}
	if !ok {
		return fmt.Errorf("file not found in MinIO bucket '%s': %s", config.BucketName, ill.FileName)
	}

	return config.DB.Create(ill).Error
}

func DeleteIllustration(id string) error {
	return config.DB.Delete(&models.Illustration{}, id).Error
}

func GetDownloadURL(fileName string, duration time.Duration) (string, error) {
	ctx := context.Background()
	reqParams := make(url.Values)

	u, err := config.MinioClient.PresignedGetObject(ctx, config.BucketName, fileName, duration, reqParams)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
