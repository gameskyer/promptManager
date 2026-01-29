package handlers

import (
	"promptmaster/backend/services"
)

// ImageHandler exposes image operations to the frontend
type ImageHandler struct {
	service *services.ImageService
}

// NewImageHandler creates a new ImageHandler
func NewImageHandler(service *services.ImageService) *ImageHandler {
	return &ImageHandler{service: service}
}

// ImageResponse represents the response structure
type ImageResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// UploadImageRequest represents an upload request
type UploadImageRequest struct {
	Data      string `json:"data"`
	PresetID  uint   `json:"preset_id"`
	VersionID uint   `json:"version_id"`
}

// UploadImage uploads a base64 encoded image
func (h *ImageHandler) UploadImage(req UploadImageRequest) ImageResponse {
	result, err := h.service.UploadImage(services.UploadImageRequest{
		Data:      req.Data,
		PresetID:  req.PresetID,
		VersionID: req.VersionID,
	})
	if err != nil {
		return ImageResponse{Success: false, Error: err.Error()}
	}
	return ImageResponse{Success: true, Data: result}
}

// GetImageByID retrieves an image by ID
func (h *ImageHandler) GetImageByID(id uint) ImageResponse {
	preview, err := h.service.GetImageByID(id)
	if err != nil {
		return ImageResponse{Success: false, Error: err.Error()}
	}
	return ImageResponse{Success: true, Data: preview}
}

// GetImagesByPreset retrieves all images for a preset
func (h *ImageHandler) GetImagesByPreset(presetID uint) ImageResponse {
	previews, err := h.service.GetImagesByPreset(presetID)
	if err != nil {
		return ImageResponse{Success: false, Error: err.Error()}
	}
	return ImageResponse{Success: true, Data: previews}
}

// DeleteImage deletes an image by ID
func (h *ImageHandler) DeleteImage(id uint) ImageResponse {
	if err := h.service.DeleteImage(id); err != nil {
		return ImageResponse{Success: false, Error: err.Error()}
	}
	return ImageResponse{Success: true}
}

// GetImageData retrieves image data as base64
func (h *ImageHandler) GetImageData(id uint) ImageResponse {
	base64Data, mimeType, err := h.service.GetImageData(id)
	if err != nil {
		return ImageResponse{Success: false, Error: err.Error()}
	}
	return ImageResponse{Success: true, Data: map[string]string{
		"data":      base64Data,
		"mime_type": mimeType,
	}}
}
