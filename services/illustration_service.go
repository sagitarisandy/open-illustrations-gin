package services

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path/filepath"
	"time"

	"open-illustrations-go/config"
	"open-illustrations-go/models"

	"github.com/minio/minio-go/v7"
)

// func UploadObject(objectName string, r io.Reader, size int64, contentType string) error {
// 	_, err := config.MinioClient.PutObject(
// 		context.Background(),
// 		config.BucketName,
// 		objectName,
// 		r,
// 		size,
// 		minio.PutObjectOptions{ContentType: contentType},
// 	)
// 	return err
// }

func UploadObject(objectName string, r io.Reader, size int64, contentType string) error {
	if contentType == "" {
		if ext := filepath.Ext(objectName); ext != "" {
			if ct := mime.TypeByExtension(ext); ct != "" {
				contentType = ct
			}
		}
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	_, err := config.MinioClient.PutObject(
		context.Background(),
		config.BucketName,
		objectName,
		r,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return err
}

func MinioObjectExists(objectName string) (bool, error) {
	_, err := config.MinioClient.StatObject(context.TODO(), config.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return false, nil
	}

	return true, nil
}

func CreateIllustrationRecord(ill *models.Illustration) error {
	return config.DB.Create(ill).Error
}

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
