package handlers

import (
	"promptmaster/backend/services"
)

// PresetWithSnapshot is re-exported for frontend
type PresetWithSnapshot = services.PresetWithSnapshot

// PresetHandler exposes preset operations to the frontend
type PresetHandler struct {
	service *services.PresetService
}

// NewPresetHandler creates a new PresetHandler
func NewPresetHandler(service *services.PresetService) *PresetHandler {
	return &PresetHandler{service: service}
}

// PresetResponse represents the response structure
type PresetResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CreatePresetRequest represents a create preset request
type CreatePresetRequest struct {
	Title      string                 `json:"title"`
	CategoryID uint                   `json:"category_id"`
	PosText    string                 `json:"pos_text"`
	NegText    string                 `json:"neg_text"`
	AtomIDs    []uint                 `json:"atom_ids"`
	Params     map[string]interface{} `json:"params"`
	Previews   []string               `json:"previews"` // base64 encoded images
}

// CreatePreset creates a new preset
func (h *PresetHandler) CreatePreset(req CreatePresetRequest) PresetResponse {
	preset, err := h.service.CreatePreset(req.Title, req.CategoryID, req.PosText, req.NegText, req.AtomIDs, req.Params, req.Previews)
	if err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	// Convert to include snapshot data
	result, err := h.service.ToPresetWithSnapshot(preset)
	if err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	return PresetResponse{Success: true, Data: result}
}

// GetPresetByID retrieves a preset by ID
func (h *PresetHandler) GetPresetByID(id uint) PresetResponse {
	preset, err := h.service.GetPresetByID(id)
	if err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	// Convert to include snapshot data
	result, err := h.service.ToPresetWithSnapshot(preset)
	if err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	return PresetResponse{Success: true, Data: result}
}

// GetPresetsRequest represents a get presets request
type GetPresetsRequest struct {
	Page           int  `json:"page"`
	PageSize       int  `json:"page_size"`
	CategoryID     uint `json:"category_id"`
	IncludeDeleted bool `json:"include_deleted"`
}

// GetPresets retrieves presets with pagination
func (h *PresetHandler) GetPresets(req GetPresetsRequest) PresetResponse {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	
	presets, total, err := h.service.GetPresets(req.Page, req.PageSize, req.CategoryID, req.IncludeDeleted)
	if err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	
	// Convert presets to include snapshot data
	var results []*services.PresetWithSnapshot
	for _, preset := range presets {
		result, err := h.service.ToPresetWithSnapshot(&preset)
		if err == nil {
			results = append(results, result)
		}
	}
	
	return PresetResponse{Success: true, Data: map[string]interface{}{
		"presets": results,
		"total":   total,
		"page":    req.Page,
	}}
}

// UpdatePresetRequest represents an update preset request
type UpdatePresetRequest struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	CategoryID uint   `json:"category_id"`
}

// UpdatePreset updates a preset
func (h *PresetHandler) UpdatePreset(req UpdatePresetRequest) PresetResponse {
	preset, err := h.service.UpdatePreset(req.ID, req.Title, req.CategoryID)
	if err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	return PresetResponse{Success: true, Data: preset}
}

// SoftDeletePreset soft deletes a preset
func (h *PresetHandler) SoftDeletePreset(id uint) PresetResponse {
	if err := h.service.SoftDeletePreset(id); err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	return PresetResponse{Success: true}
}

// RestorePreset restores a soft-deleted preset
func (h *PresetHandler) RestorePreset(id uint) PresetResponse {
	if err := h.service.RestorePreset(id); err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	return PresetResponse{Success: true}
}

// BuildPromptRequest represents a build prompt request
type BuildPromptRequest struct {
	AtomIDs []uint `json:"atom_ids"`
}

// BuildPrompt builds the prompt text from atom IDs
func (h *PresetHandler) BuildPrompt(req BuildPromptRequest) PresetResponse {
	text, ids, err := h.service.BuildPromptText(req.AtomIDs)
	if err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	
	return PresetResponse{Success: true, Data: map[string]interface{}{
		"text":     text,
		"atom_ids": ids,
	}}
}

// GetCurrentWorkState retrieves the current working state for a preset
func (h *PresetHandler) GetCurrentWorkState(presetID uint) PresetResponse {
	state, versionNum, err := h.service.GetCurrentWorkState(presetID)
	if err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	
	return PresetResponse{Success: true, Data: map[string]interface{}{
		"state":        state,
		"version_num":  versionNum,
	}}
}

// ForkPresetRequest represents a fork preset request
type ForkPresetRequest struct {
	PresetID   uint   `json:"preset_id"`
	VersionNum int    `json:"version_num"`
	NewTitle   string `json:"new_title"`
}

// ForkPreset creates a new preset based on an existing version
func (h *PresetHandler) ForkPreset(req ForkPresetRequest) PresetResponse {
	// 获取原预设的分类ID
	sourcePreset, err := h.service.GetPresetByID(req.PresetID)
	if err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	
	preset, err := h.service.ForkPreset(req.PresetID, req.VersionNum, req.NewTitle, sourcePreset.CategoryID)
	if err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	return PresetResponse{Success: true, Data: preset}
}

// CleanupOldVersionsRequest represents a cleanup request
type CleanupOldVersionsRequest struct {
	PresetID  uint `json:"preset_id"`
	KeepCount int  `json:"keep_count"`
}

// CleanupOldVersions removes old versions based on retention policy
func (h *PresetHandler) CleanupOldVersions(req CleanupOldVersionsRequest) PresetResponse {
	if err := h.service.CleanupOldVersions(req.PresetID, req.KeepCount); err != nil {
		return PresetResponse{Success: false, Error: err.Error()}
	}
	return PresetResponse{Success: true}
}
