package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client
var BucketName = os.Getenv("MINIO_BUCKET")

func InitMinio() {
	_ = godotenv.Load()

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ROOT_USER")
	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	useSSL, _ := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Error connect MinIO: %v", err)
	}
	MinioClient = client

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, BucketName)
	if err != nil {
		log.Fatalf("Error check bucket: %v", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{Region: "us-east-1"}); err != nil {
			log.Fatalf("Error create bucket: %v", err)
		}
		log.Printf("Bucket %s created!", BucketName)
	} else {
		log.Printf("Bucket %s already exists", BucketName)
	}
}
