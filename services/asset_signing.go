package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"open-illustrations-go/config"

	"github.com/minio/minio-go/v7"
)

func assetSecret() ([]byte, error) {
	sec := os.Getenv("ASSET_SIGNING_SECRET")
	if sec == "" {
		return nil, errors.New("ASSET_SIGNING_SECRET not set")
	}
	return []byte(sec), nil
}

// GenerateAssetToken creates a short-lived signed token for a storageKey.
func GenerateAssetToken(storageKey string, ttl time.Duration) (string, error) {
	secret, err := assetSecret()
	if err != nil {
		return "", err
	}
	exp := time.Now().Add(ttl).Unix()
	payload := storageKey + "." + strconv.FormatInt(exp, 10)
	m := hmac.New(sha256.New, secret)
	m.Write([]byte(payload))
	sig := m.Sum(nil)
	raw := storageKey + "|" + strconv.FormatInt(exp, 10) + "|" + base64.RawURLEncoding.EncodeToString(sig)
	return base64.RawURLEncoding.EncodeToString([]byte(raw)), nil
}

// ParseAndValidateAssetToken validates token and returns storageKey if valid.
func ParseAndValidateAssetToken(token string) (string, error) {
	secret, err := assetSecret()
	if err != nil {
		return "", err
	}
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", errors.New("invalid token encoding")
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 3 {
		return "", errors.New("invalid token parts")
	}
	storageKey := parts[0]
	expStr := parts[1]
	sigStr := parts[2]

	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return "", errors.New("invalid exp")
	}
	if time.Now().Unix() > exp {
		return "", errors.New("expired")
	}

	payload := storageKey + "." + expStr
	m := hmac.New(sha256.New, secret)
	m.Write([]byte(payload))
	expected := m.Sum(nil)
	got, err := base64.RawURLEncoding.DecodeString(sigStr)
	if err != nil {
		return "", errors.New("invalid sig encoding")
	}
	if !hmac.Equal(expected, got) {
		return "", errors.New("signature mismatch")
	}
	return storageKey, nil
}

// GetObjectStream returns a readable MinIO object stream with its content-type.
func GetObjectStream(storageKey string) (*minio.Object, string, io.ReadSeeker, error) {
	obj, err := config.MinioClient.GetObject(context.Background(), config.BucketName, storageKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", nil, err
	}
	st, err := obj.Stat()
	if err != nil {
		obj.Close()
		return nil, "", nil, err
	}
	return obj, st.ContentType, obj, nil
}
