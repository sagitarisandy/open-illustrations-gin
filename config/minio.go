package config

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client
var BucketName = "illustrations"

func InitMinio() {
	endpoint := "VPS_IP"
	accessKey := "ACCESS_KEY"
	secretKey := "SECRET_KEY"
	useSSL := false

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
