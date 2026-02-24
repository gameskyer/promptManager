package handlers

import (
	"promptmaster/backend/services"
)

// BatchHandler exposes batch operations to the frontend
type BatchHandler struct {
	service *services.BatchService
}

// NewBatchHandler creates a new BatchHandler
func NewBatchHandler(service *services.BatchService) *BatchHandler {
	return &BatchHandler{service: service}
}

// BatchResponse represents the response structure
type BatchResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// BatchMoveCategoryRequest represents a batch move request
type BatchMoveCategoryRequest struct {
	AtomIDs    []uint `json:"atom_ids"`
	CategoryID uint   `json:"category_id"`
}

// BatchMoveCategory moves multiple atoms to a new category
func (h *BatchHandler) BatchMoveCategory(req BatchMoveCategoryRequest) BatchResponse {
	count, err := h.service.BatchMoveCategory(req.AtomIDs, req.CategoryID)
	if err != nil {
		return BatchResponse{Success: false, Error: err.Error()}
	}
	return BatchResponse{
		Success: true,
		Data: map[string]interface{}{
			"moved_count": count,
		},
	}
}

// BatchUpdateTypeRequest represents a batch type update request
type BatchUpdateTypeRequest struct {
	AtomIDs []uint `json:"atom_ids"`
	Type    string `json:"type"` // "Positive" or "Negative"
}

// BatchUpdateType updates the type of multiple atoms
func (h *BatchHandler) BatchUpdateType(req BatchUpdateTypeRequest) BatchResponse {
	count, err := h.service.BatchUpdateType(req.AtomIDs, req.Type)
	if err != nil {
		return BatchResponse{Success: false, Error: err.Error()}
	}
	return BatchResponse{
		Success: true,
		Data: map[string]interface{}{
			"updated_count": count,
		},
	}
}

// BatchDeleteRequest represents a batch delete request
type BatchDeleteRequest struct {
	AtomIDs []uint `json:"atom_ids"`
}

// BatchDelete soft deletes multiple atoms
func (h *BatchHandler) BatchDelete(req BatchDeleteRequest) BatchResponse {
	count, err := h.service.BatchDelete(req.AtomIDs)
	if err != nil {
		return BatchResponse{Success: false, Error: err.Error()}
	}
	return BatchResponse{
		Success: true,
		Data: map[string]interface{}{
			"deleted_count": count,
		},
	}
}

// BatchAddSynonymsRequest represents a batch add synonyms request
type BatchAddSynonymsRequest struct {
	AtomIDs  []uint   `json:"atom_ids"`
	Synonyms []string `json:"synonyms"`
}

// BatchAddSynonyms adds synonyms to multiple atoms
func (h *BatchHandler) BatchAddSynonyms(req BatchAddSynonymsRequest) BatchResponse {
	count, err := h.service.BatchAddSynonyms(req.AtomIDs, req.Synonyms)
	if err != nil {
		return BatchResponse{Success: false, Error: err.Error()}
	}
	return BatchResponse{
		Success: true,
		Data: map[string]interface{}{
			"updated_count": count,
		},
	}
}

// BatchClearCategory clears the category of multiple atoms
func (h *BatchHandler) BatchClearCategory(req BatchDeleteRequest) BatchResponse {
	count, err := h.service.BatchClearCategory(req.AtomIDs)
	if err != nil {
		return BatchResponse{Success: false, Error: err.Error()}
	}
	return BatchResponse{
		Success: true,
		Data: map[string]interface{}{
			"cleared_count": count,
		},
	}
}
