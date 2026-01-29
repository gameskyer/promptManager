package services

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"promptmaster/backend/config"
	"promptmaster/backend/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ImageService handles image upload and management
type ImageService struct {
	db *gorm.DB
}

// NewImageService creates a new ImageService
func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{db: db}
}

// UploadImageRequest represents an image upload request
type UploadImageRequest struct {
	Data      string `json:"data"`       // base64 encoded image data
	PresetID  uint   `json:"preset_id"`  // optional preset association
	VersionID uint   `json:"version_id"` // optional version association
}

// UploadResult represents the upload result
type UploadResult struct {
	ID       uint   `json:"id"`
	FilePath string `json:"file_path"`
	URL      string `json:"url"`
}

// Allowed image mime types
var allowedMimeTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// UploadImage saves a base64 encoded image to disk
func (s *ImageService) UploadImage(req UploadImageRequest) (*UploadResult, error) {
	// Parse base64 data
	data, mimeType, err := parseBase64Image(req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image: %w", err)
	}

	// Validate mime type
	ext, ok := allowedMimeTypes[mimeType]
	if !ok {
		return nil, fmt.Errorf("unsupported image type: %s", mimeType)
	}

	// Generate unique filename
	filename := uuid.New().String() + ext
	filePath := filepath.Join(config.ImageDir, filename)

	// Save file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	// Save to database
	preview := &models.Preview{
		PresetID:  req.PresetID,
		VersionID: func() *uint { if req.VersionID > 0 { return &req.VersionID }; return nil }(),
		FilePath:  filePath,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(preview).Error; err != nil {
		// Clean up file on DB error
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save to database: %w", err)
	}

	return &UploadResult{
		ID:       preview.ID,
		FilePath: filePath,
		URL:      "file://" + filePath,
	}, nil
}

// UploadImageFromPath copies an image from a path to the image directory
func (s *ImageService) UploadImageFromPath(sourcePath string, presetID uint, versionID uint) (*UploadResult, error) {
	// Open source file
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Detect mime type
	buffer := make([]byte, 512)
	n, err := srcFile.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	mimeType := http.DetectContentType(buffer[:n])

	// Validate mime type
	ext, ok := allowedMimeTypes[mimeType]
	if !ok {
		return nil, fmt.Errorf("unsupported image type: %s", mimeType)
	}

	// Reset file pointer
	srcFile.Seek(0, 0)

	// Generate unique filename
	filename := uuid.New().String() + ext
	filePath := filepath.Join(config.ImageDir, filename)

	// Create destination file
	dstFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy file content
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Save to database
	preview := &models.Preview{
		PresetID:  presetID,
		VersionID: func() *uint { if versionID > 0 { return &versionID }; return nil }(),
		FilePath:  filePath,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(preview).Error; err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save to database: %w", err)
	}

	return &UploadResult{
		ID:       preview.ID,
		FilePath: filePath,
		URL:      "file://" + filePath,
	}, nil
}

// GetImageByID retrieves an image by ID
func (s *ImageService) GetImageByID(id uint) (*models.Preview, error) {
	var preview models.Preview
	if err := s.db.First(&preview, id).Error; err != nil {
		return nil, err
	}
	return &preview, nil
}

// GetImagesByPreset retrieves all images for a preset
func (s *ImageService) GetImagesByPreset(presetID uint) ([]models.Preview, error) {
	var previews []models.Preview
	if err := s.db.Where("preset_id = ?", presetID).Find(&previews).Error; err != nil {
		return nil, err
	}
	return previews, nil
}

// DeleteImage deletes an image by ID
func (s *ImageService) DeleteImage(id uint) error {
	var preview models.Preview
	if err := s.db.First(&preview, id).Error; err != nil {
		return err
	}

	// Delete file
	os.Remove(preview.FilePath)

	// Delete from database
	return s.db.Delete(&preview).Error
}

// GetImageData reads image data as base64
func (s *ImageService) GetImageData(id uint) (string, string, error) {
	preview, err := s.GetImageByID(id)
	if err != nil {
		return "", "", err
	}

	data, err := os.ReadFile(preview.FilePath)
	if err != nil {
		return "", "", err
	}

	// Detect mime type
	mimeType := http.DetectContentType(data)
	base64Data := base64.StdEncoding.EncodeToString(data)

	return base64Data, mimeType, nil
}

// parseBase64Image parses base64 image data
func parseBase64Image(dataURI string) ([]byte, string, error) {
	// Check if it's a data URI
	if strings.HasPrefix(dataURI, "data:") {
		// Parse data URI
		commaIndex := strings.Index(dataURI, ",")
		if commaIndex == -1 {
			return nil, "", fmt.Errorf("invalid data URI")
		}

		// Extract mime type
		meta := dataURI[5:commaIndex]
		parts := strings.Split(meta, ";")
		mimeType := parts[0]

		// Decode base64
		data, err := base64.StdEncoding.DecodeString(dataURI[commaIndex+1:])
		if err != nil {
			return nil, "", err
		}

		return data, mimeType, nil
	}

	// Plain base64, try to decode
	data, err := base64.StdEncoding.DecodeString(dataURI)
	if err != nil {
		return nil, "", err
	}

	// Detect mime type from data
	mimeType := http.DetectContentType(data)
	return data, mimeType, nil
}
