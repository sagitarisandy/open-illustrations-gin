package services

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"open-illustrations-go/config"
	"open-illustrations-go/models"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
	result := config.DB.
		Preload("CategoryRef").
		Preload("PackRef").
		Preload("StyleRef").
		Where("deleted_at IS NULL").
		Find(&illustrations)
	return illustrations, result.Error
}

func GetIllustrationsByCategory(categoryID string) ([]models.Illustration, error) {
	var illustrations []models.Illustration
	result := config.DB.
		Preload("CategoryRef").
		Preload("PackRef").
		Preload("StyleRef").
		Where("deleted_at IS NULL AND category_id = ?", categoryID).
		Find(&illustrations)
	return illustrations, result.Error
}

func GetIllustrationsByStyle(styleID string) ([]models.Illustration, error) {
	var illustrations []models.Illustration
	result := config.DB.
		Preload("CategoryRef").
		Preload("PackRef").
		Preload("StyleRef").
		Where("deleted_at IS NULL AND style_id = ?", styleID).
		Find(&illustrations)
	return illustrations, result.Error
}

func GetIllustrationsByPack(packID string) ([]models.Illustration, error) {
	var illustrations []models.Illustration
	result := config.DB.
		Preload("CategoryRef").
		Preload("PackRef").
		Preload("StyleRef").
		Where("deleted_at IS NULL AND pack_id = ?", packID).
		Find(&illustrations)
	return illustrations, result.Error
}

func GetIllustration(id string) (*models.Illustration, error) {
	var illustration models.Illustration
	result := config.DB.
		Preload("CategoryRef").
		Preload("PackRef").
		Preload("StyleRef").
		First(&illustration, id)
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
	ok, err := minioObjectExists(ill.StorageKey)
	if err != nil {
		return fmt.Errorf("minio check failed: %w", err)
	}
	if !ok {
		return fmt.Errorf("file not found in MinIO bucket '%s': %s", config.BucketName, ill.StorageKey)
	}

	return config.DB.Create(ill).Error
}

func DeleteIllustration(id string) error {
	return config.DB.Delete(&models.Illustration{}, id).Error
}

func GetDownloadURL(storageKey string, duration time.Duration) (string, error) {
	ctx := context.Background()
	reqParams := make(url.Values)

	cli := config.MinioClient // default: endpoint internal (mini:9000)

	// Jika MINIO_PUBLIC_BASE_URL diset (mis. http://localhost:9000), buat client khusus presign
	if base := os.Getenv("MINIO_PUBLIC_BASE_URL"); base != "" {
		if pub, err := url.Parse(base); err == nil && pub.Host != "" {
			pubCli, err := minio.New(pub.Host, &minio.Options{
				Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ROOT_USER"), os.Getenv("MINIO_ROOT_PASSWORD"), ""),
				Secure: pub.Scheme == "https",
				Region: "us-east-1",
			})
			if err == nil {
				cli = pubCli
			}
		}
	}

	u, err := cli.PresignedGetObject(ctx, config.BucketName, storageKey, duration, reqParams)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// PresignTTL returns the server-enforced TTL for presigned URLs.
// It ignores any client-provided values and reads from PRESIGN_TTL_SECONDS env var.
// Clamped to [60, 3600] seconds; defaults to 600 if unset/invalid.
func PresignTTL() time.Duration {
	v := os.Getenv("PRESIGN_TTL_SECONDS")
	n := 600
	if v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			n = parsed
		}
	}
	if n < 60 {
		n = 60
	}
	if n > 3600 {
		n = 3600
	}
	return time.Duration(n) * time.Second
}
