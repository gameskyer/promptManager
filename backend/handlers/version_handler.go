package handlers

import (
	"promptmaster/backend/services"
)

// VersionHandler exposes version control operations to the frontend
type VersionHandler struct {
	service *services.VersionService
}

// NewVersionHandler creates a new VersionHandler
func NewVersionHandler(service *services.VersionService) *VersionHandler {
	return &VersionHandler{service: service}
}

// VersionResponse represents the response structure
type VersionResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CreateVersionRequest represents a create version request
type CreateVersionRequest struct {
	PresetID      uint                   `json:"preset_id"`
	PosText       string                 `json:"pos_text"`
	NegText       string                 `json:"neg_text"`
	AtomIDs       []uint                 `json:"atom_ids"`
	Params        map[string]interface{} `json:"params"`
	PreviewPaths  []string               `json:"preview_paths"`
	ThumbnailPath string                 `json:"thumbnail_path"`
}

// CreateVersion creates a new version
func (h *VersionHandler) CreateVersion(req CreateVersionRequest) VersionResponse {
	input := services.VersionInput{
		PresetID:      req.PresetID,
		PosText:       req.PosText,
		NegText:       req.NegText,
		AtomIDs:       req.AtomIDs,
		Params:        req.Params,
		PreviewPaths:  req.PreviewPaths,
		ThumbnailPath: req.ThumbnailPath,
	}
	
	version, err := h.service.CreateVersion(input)
	if err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true, Data: version}
}

// GetVersion retrieves a specific version
func (h *VersionHandler) GetVersion(presetID uint, versionNum int) VersionResponse {
	version, err := h.service.GetVersion(presetID, versionNum)
	if err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true, Data: version}
}

// GetVersionHistory retrieves version history for a preset
func (h *VersionHandler) GetVersionHistory(presetID uint, limit int) VersionResponse {
	if limit <= 0 {
		limit = 20
	}
	
	versions, err := h.service.GetVersionHistory(presetID, limit)
	if err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true, Data: versions}
}

// GetLatestVersions retrieves the latest versions
func (h *VersionHandler) GetLatestVersions(presetID uint, count int) VersionResponse {
	if count <= 0 {
		count = 5
	}
	
	versions, err := h.service.GetLatestVersions(presetID, count)
	if err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true, Data: versions}
}

// StarVersionRequest represents a star version request
type StarVersionRequest struct {
	VersionID uint `json:"version_id"`
	Starred   bool `json:"starred"`
}

// StarVersion toggles the star status of a version
func (h *VersionHandler) StarVersion(req StarVersionRequest) VersionResponse {
	if err := h.service.StarVersion(req.VersionID, req.Starred); err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true}
}

// RollbackToVersion rolls back to a specific version
func (h *VersionHandler) RollbackToVersion(presetID uint, targetVersionNum int) VersionResponse {
	version, err := h.service.RollbackToVersion(presetID, targetVersionNum)
	if err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true, Data: version}
}

// CompareVersions compares two versions
func (h *VersionHandler) CompareVersions(presetID uint, v1, v2 int) VersionResponse {
	diff, err := h.service.CompareVersions(presetID, v1, v2)
	if err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true, Data: diff}
}

// DeleteVersion deletes a specific version
func (h *VersionHandler) DeleteVersion(versionID uint) VersionResponse {
	if err := h.service.DeleteVersion(versionID); err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true}
}

// GetStarredVersions returns all starred versions
func (h *VersionHandler) GetStarredVersions(presetID uint) VersionResponse {
	versions, err := h.service.GetStarredVersions(presetID)
	if err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true, Data: versions}
}

// GetVersionDiffStats gets diff statistics for a version
func (h *VersionHandler) GetVersionDiffStats(presetID uint, versionNum int) VersionResponse {
	version, err := h.service.GetVersion(presetID, versionNum)
	if err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true, Data: map[string]string{
		"diff_stats": version.DiffStats,
	}}
}

// UpdateVersionPreviewRequest represents a request to update version preview
// This updates the current version without creating a new one
type UpdateVersionPreviewRequest struct {
	PresetID      uint     `json:"preset_id"`
	ThumbnailPath string   `json:"thumbnail_path"`
	PreviewPaths  []string `json:"preview_paths"`
}

// UpdateVersionPreview updates only the preview images of the current version
// Preview changes do not create new versions
func (h *VersionHandler) UpdateVersionPreview(req UpdateVersionPreviewRequest) VersionResponse {
	version, err := h.service.UpdateVersionPreview(req.PresetID, req.ThumbnailPath, req.PreviewPaths)
	if err != nil {
		return VersionResponse{Success: false, Error: err.Error()}
	}
	return VersionResponse{Success: true, Data: version}
}
