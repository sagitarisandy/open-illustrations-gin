package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"open-illustrations-go/models"
	"open-illustrations-go/services"

	"github.com/gin-gonic/gin"
)

type CreateIllustrationDTO struct {
	Title      string `json:"title" binding:"required"`
	StyleID    *uint  `json:"style_id"`
	CategoryID *uint  `json:"category_id"`
	PackID     *uint  `json:"pack_id"`
	FileName   string `json:"file_name" binding:"required"`
	StorageKey string `json:"storage_key"`
}

func GetIllustrations(c *gin.Context) {
	includePresign := c.Query("include_presign") == "1" // for premium items only

	data, err := services.GetIllustrations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Presign TTL (server-enforced) computed only if needed
	var exp time.Duration
	if includePresign {
		exp = services.PresignTTL()
	}

	out := make([]gin.H, 0, len(data))
	for _, ill := range data {
		item := gin.H{
			"id":          ill.ID,
			"title":       ill.Title,
			"style_id":    ill.StyleID,
			"category_id": ill.CategoryID,
			"pack_id":     ill.PackID,
			"file_name":   ill.FileName,
			"is_premium":  ill.IsPremium,
			"created_at":  ill.CreatedAt,
			"updated_at":  ill.UpdatedAt,
		}

		// Decide a single image_url
		if ill.IsPremium {
			var url string
			if includePresign {
				if u, err := services.GetDownloadURL(ill.StorageKey, exp); err == nil {
					url = u
				}
			}
			if url == "" { // fallback to backend-signed token path
				if tok, err := services.GenerateAssetToken(ill.StorageKey, 15*time.Minute); err == nil {
					url = "/api/v1/i/" + tok
				}
			}
			if url == "" { // last resort, public by id
				url = fmt.Sprintf("/api/v1/illustrations/%d/public", ill.ID)
			}
			item["image_url"] = url
		} else {
			item["image_url"] = fmt.Sprintf("/api/v1/illustrations/%d/public", ill.ID)
		}

		out = append(out, item)
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// GetIllustrationsByCategory handles GET /api/v1/categories/:id/illustrations
func GetIllustrationsByCategory(c *gin.Context) {
	includePresign := c.Query("include_presign") == "1"
	id := c.Param("id")

	data, err := services.GetIllustrationsByCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var exp time.Duration
	if includePresign {
		exp = services.PresignTTL()
	}
	out := make([]gin.H, 0, len(data))
	for _, ill := range data {
		item := gin.H{
			"id": ill.ID, "title": ill.Title, "style_id": ill.StyleID,
			"category_id": ill.CategoryID, "pack_id": ill.PackID,
			"file_name": ill.FileName, "is_premium": ill.IsPremium,
			"created_at": ill.CreatedAt, "updated_at": ill.UpdatedAt,
		}
		if ill.IsPremium {
			var url string
			if includePresign {
				if u, err := services.GetDownloadURL(ill.StorageKey, exp); err == nil {
					url = u
				}
			}
			if url == "" {
				if tok, err := services.GenerateAssetToken(ill.StorageKey, 15*time.Minute); err == nil {
					url = "/api/v1/i/" + tok
				}
			}
			if url == "" {
				url = fmt.Sprintf("/api/v1/illustrations/%d/public", ill.ID)
			}
			item["image_url"] = url
		} else {
			item["image_url"] = fmt.Sprintf("/api/v1/illustrations/%d/public", ill.ID)
		}
		out = append(out, item)
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// GetIllustrationsByStyle handles GET /api/v1/styles/:id/illustrations
func GetIllustrationsByStyle(c *gin.Context) {
	includePresign := c.Query("include_presign") == "1"
	id := c.Param("id")
	data, err := services.GetIllustrationsByStyle(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var exp time.Duration
	if includePresign {
		exp = services.PresignTTL()
	}
	out := make([]gin.H, 0, len(data))
	for _, ill := range data {
		item := gin.H{"id": ill.ID, "title": ill.Title, "style_id": ill.StyleID, "category_id": ill.CategoryID, "pack_id": ill.PackID, "file_name": ill.FileName, "is_premium": ill.IsPremium, "created_at": ill.CreatedAt, "updated_at": ill.UpdatedAt}
		if ill.IsPremium {
			var url string
			if includePresign {
				if u, err := services.GetDownloadURL(ill.StorageKey, exp); err == nil {
					url = u
				}
			}
			if url == "" {
				if tok, err := services.GenerateAssetToken(ill.StorageKey, 15*time.Minute); err == nil {
					url = "/api/v1/i/" + tok
				}
			}
			if url == "" {
				url = fmt.Sprintf("/api/v1/illustrations/%d/public", ill.ID)
			}
			item["image_url"] = url
		} else {
			item["image_url"] = fmt.Sprintf("/api/v1/illustrations/%d/public", ill.ID)
		}
		out = append(out, item)
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// GetIllustrationsByPack handles GET /api/v1/packs/:id/illustrations
func GetIllustrationsByPack(c *gin.Context) {
	includePresign := c.Query("include_presign") == "1"
	id := c.Param("id")
	data, err := services.GetIllustrationsByPack(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var exp time.Duration
	if includePresign {
		exp = services.PresignTTL()
	}
	out := make([]gin.H, 0, len(data))
	for _, ill := range data {
		item := gin.H{"id": ill.ID, "title": ill.Title, "style_id": ill.StyleID, "category_id": ill.CategoryID, "pack_id": ill.PackID, "file_name": ill.FileName, "is_premium": ill.IsPremium, "created_at": ill.CreatedAt, "updated_at": ill.UpdatedAt}
		if ill.IsPremium {
			var url string
			if includePresign {
				if u, err := services.GetDownloadURL(ill.StorageKey, exp); err == nil {
					url = u
				}
			}
			if url == "" {
				if tok, err := services.GenerateAssetToken(ill.StorageKey, 15*time.Minute); err == nil {
					url = "/api/v1/i/" + tok
				}
			}
			if url == "" {
				url = fmt.Sprintf("/api/v1/illustrations/%d/public", ill.ID)
			}
			item["image_url"] = url
		} else {
			item["image_url"] = fmt.Sprintf("/api/v1/illustrations/%d/public", ill.ID)
		}
		out = append(out, item)
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func GetIllustration(c *gin.Context) {
	id := c.Param("id")
	ill, err := services.GetIllustration(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	// Build single image_url based on premium flag
	payload := gin.H{
		"id":          ill.ID,
		"title":       ill.Title,
		"style_id":    ill.StyleID,
		"category_id": ill.CategoryID,
		"pack_id":     ill.PackID,
		"file_name":   ill.FileName,
		"is_premium":  ill.IsPremium,
		"created_at":  ill.CreatedAt,
		"updated_at":  ill.UpdatedAt,
	}
	// Decide image_url
	var url string
	if ill.IsPremium {
		if c.Query("include_presign") == "1" {
			exp := services.PresignTTL()
			if u, err := services.GetDownloadURL(ill.StorageKey, exp); err == nil {
				url = u
			}
		}
		if url == "" { // fallback to backend token
			if tok, err := services.GenerateAssetToken(ill.StorageKey, 15*time.Minute); err == nil {
				url = "/api/v1/i/" + tok
			}
		}
		if url == "" { // last resort
			url = fmt.Sprintf("/api/v1/illustrations/%s/public", id)
		}
	} else {
		url = fmt.Sprintf("/api/v1/illustrations/%s/public", id)
	}
	payload["image_url"] = url
	c.JSON(http.StatusOK, gin.H{"data": payload})
}

// GetIllustrationFileURL returns a short-lived presigned URL for a given storage key
// Route: GET /api/v1/illustrations/file/:key
func GetIllustrationFileURL(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing key"})
		return
	}

	// Optional: ?expires=seconds (clamp 60..3600)
	// Enforce server-side policy for presign TTL
	exp := services.PresignTTL()

	u, err := services.GetDownloadURL(key, exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// also include backend signed proxy path as a fallback that doesn't expose storage details
	tok, _ := services.GenerateAssetToken(key, 15*time.Minute)
	c.JSON(http.StatusOK, gin.H{
		"url":        u,
		"expires_in": int(exp.Seconds()),
		"signed_url": "/api/v1/i/" + tok,
	})
}

// GetIllustrationFileURLByID returns a short-lived presigned URL for an illustration by its numeric ID
// Route: GET /api/v1/illustrations/:id/file
func GetIllustrationFileURLByID(c *gin.Context) {
	id := c.Param("id")
	// lookup illustration to get its storage key
	ill, err := services.GetIllustration(id)
	if err != nil || ill == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "illustration not found"})
		return
	}

	exp := services.PresignTTL()

	u, err := services.GetDownloadURL(ill.StorageKey, exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tok, _ := services.GenerateAssetToken(ill.StorageKey, 15*time.Minute)
	c.JSON(http.StatusOK, gin.H{
		"url":        u,
		"expires_in": int(exp.Seconds()),
		"signed_url": "/api/v1/i/" + tok,
	})
}

// Removed explicit GetIllustrationURL in favor of signed URL embedded responses

// processUpload handles multipart form upload: fields => file, title, category, file_name(optional)
func processUpload(c *gin.Context) {

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form field 'file' is required"})
		return
	}
	title := c.PostForm("title")
	styleIDStr := c.PostForm("style_id")
	catIDStr := c.PostForm("category_id")
	packIDStr := c.PostForm("pack_id")
	objectName := c.PostForm("file_name")
	if objectName == "" {
		objectName = fh.Filename
	}
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form field 'title' is required"})
		return
	}

	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer f.Close()

	// Validate SVG only (by extension + light content-type check)
	if !isSVGFile(fh) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .svg files are allowed"})
		return
	}

	// Generate unique storage key (random hex) while preserving original filename separately
	storageKey := generateStorageKey(objectName)

	// Jika sudah ada nama object yang sama -> tolak (hindari duplikat)
	exists, err := services.MinioObjectExists(storageKey)
	if err != nil {
		log.Println("check minio exists err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "minio check failed"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "file already exists in bucket", "file_name": objectName})
		return
	}

	// upload ke MinIO
	if err := services.UploadObject(storageKey, f, fh.Size, fh.Header.Get("Content-Type")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload to storage", "detail": err.Error()})
		return
	}

	var catIDPtr *uint
	var styleIDPtr *uint
	var packIDPtr *uint
	if catIDStr != "" {
		if v, err := strconv.ParseUint(catIDStr, 10, 64); err == nil {
			vv := uint(v)
			catIDPtr = &vv
		}
	}
	if styleIDStr != "" {
		if v, err := strconv.ParseUint(styleIDStr, 10, 64); err == nil {
			u := uint(v)
			styleIDPtr = &u
		}
	}
	if packIDStr != "" {
		if v, err := strconv.ParseUint(packIDStr, 10, 64); err == nil {
			vv := uint(v)
			packIDPtr = &vv
		}
	}
	rec := models.Illustration{
		Title:      title,
		StyleID:    styleIDPtr,
		CategoryID: catIDPtr,
		PackID:     packIDPtr,
		FileName:   objectName,
		StorageKey: storageKey,
	}
	if err := services.CreateIllustration(&rec); err != nil {
		log.Println("db insert err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": rec})
}

// Deprecated path: POST /illustrations/upload (still works). Prefer using POST /illustrations with multipart form-data.
func UploadIllustration(c *gin.Context) {
	processUpload(c)
}

func CreateIllustration(c *gin.Context) {
	// If client sent multipart form (file upload), reuse upload logic here so
	// people can just POST /illustrations with form-data.
	if ct := c.GetHeader("Content-Type"); strings.Contains(ct, "multipart/form-data") {
		processUpload(c)
		return
	}

	var dto CreateIllustrationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload (make sure Content-Type: application/json)"})
		return
	}
	// c.JSON(http.StatusCreated, gin.H{"data": input})

	storageKey := dto.StorageKey
	if storageKey == "" {
		// For JSON-based creation we expect the storage_key (object already uploaded via another service)
		c.JSON(http.StatusBadRequest, gin.H{"error": "storage_key is required when creating via JSON"})
		return
	}
	input := models.Illustration{
		Title:      dto.Title,
		StyleID:    dto.StyleID,
		CategoryID: dto.CategoryID,
		PackID:     dto.PackID,
		FileName:   dto.FileName,
		StorageKey: storageKey,
	}

	if err := services.CreateIllustration(&input); err != nil {
		log.Println("CreateIllustration DB/MINIO err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func DeleteIllustration(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteIllustration(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete illustration"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func Download(c *gin.Context) {
	id := c.Param("id")
	ill, err := services.GetIllustration(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	url, err := services.GetDownloadURL(ill.StorageKey, time.Hour*1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate link"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"download_url": url})
}

// StreamSigned serves image via backend using signed token path: /api/v1/i/:token
func StreamSigned(c *gin.Context) {
	token := c.Param("token")
	storageKey, err := services.ParseAndValidateAssetToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}
	obj, ct, reader, err := services.GetObjectStream(storageKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "object not found"})
		return
	}
	defer obj.Close()

	if ct == "" || !strings.Contains(ct, "svg") {
		ct = "image/svg+xml"
	}

	sum := sha256.Sum256([]byte(storageKey))
	etag := base64.RawURLEncoding.EncodeToString(sum[:8])
	if inm := c.GetHeader("If-None-Match"); inm != "" && inm == etag {
		c.Status(http.StatusNotModified)
		return
	}
	c.Header("Content-Type", ct)
	c.Header("Cache-Control", "public, max-age=900")
	c.Header("ETag", etag)
	c.Header("Content-Disposition", "inline; filename=\""+storageKey+"\"")
	c.Header("Content-Security-Policy", "default-src 'none'; img-src 'self'; style-src 'unsafe-inline'")
	_, _ = io.Copy(c.Writer, reader)
}

// StreamPublic serves non-premium images publicly by illustration ID: /api/v1/illustrations/:id/public
func StreamPublic(c *gin.Context) {
	id := c.Param("id")
	ill, err := services.GetIllustration(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if ill.IsPremium {
		c.JSON(http.StatusForbidden, gin.H{"error": "premium content is not publicly accessible"})
		return
	}
	obj, ct, reader, err := services.GetObjectStream(ill.StorageKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "object not found"})
		return
	}
	defer obj.Close()

	if ct == "" || !strings.Contains(ct, "svg") {
		ct = "image/svg+xml"
	}
	sum := sha256.Sum256([]byte(ill.StorageKey))
	etag := base64.RawURLEncoding.EncodeToString(sum[:8])
	if inm := c.GetHeader("If-None-Match"); inm != "" && inm == etag {
		c.Status(http.StatusNotModified)
		return
	}
	c.Header("Content-Type", ct)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Header("ETag", etag)
	c.Header("Content-Disposition", "inline; filename=\""+ill.FileName+"\"")
	c.Header("Content-Security-Policy", "default-src 'none'; img-src 'self'; style-src 'unsafe-inline'")
	_, _ = io.Copy(c.Writer, reader)
}

// --- helpers ---
func isSVGFile(fh *multipart.FileHeader) bool {
	name := strings.ToLower(fh.Filename)
	if !strings.HasSuffix(name, ".svg") {
		return false
	}
	f, err := fh.Open()
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	snippet := strings.ToLower(string(buf[:n]))
	return strings.Contains(snippet, "<svg")
}

func generateStorageKey(original string) string {
	// random 8 bytes -> 16 hex chars + original ext
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return original // fallback
	}
	ext := filepath.Ext(original)
	if ext == "" {
		ext = ".svg" // default
	}
	return fmt.Sprintf("%s-%s%s", time.Now().Format("20060102"), hex.EncodeToString(b), ext)
}
