package config

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/joho/godotenv"
)

var MinioClient *minio.Client
var BucketName = "illustrations"

func InitMinio() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file is not found")
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ROOT_USER")
	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV2(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Error connect MinIO: %v", err)
	}

	MinioClient = client

	// ensure bucket is available
	ctx := context.Background()
	exists, errBucket := client.BucketExists(ctx, BucketName)
	if errBucket != nil {
		log.Fatalf("Error check bucket: %v", errBucket)
	}

	if !exists {
		err = client.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			log.Fatalf("Error create bucket: %v", err)
		}
		log.Printf("Bucket %s created!", BucketName)
	} else {
		log.Printf("Bucket %s already exists", BucketName)
	}
}
