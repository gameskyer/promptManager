package handlers

import (
	"encoding/json"
	"promptmaster/backend/services"
)

// AtomHandler exposes atom operations to the frontend
type AtomHandler struct {
	service *services.AtomService
}

// NewAtomHandler creates a new AtomHandler
func NewAtomHandler(service *services.AtomService) *AtomHandler {
	return &AtomHandler{service: service}
}

// AtomResponse represents the response structure
type AtomResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CreateAtomRequest represents a create atom request
type CreateAtomRequest struct {
	Value      string   `json:"value"`
	Label      string   `json:"label"`
	Type       string   `json:"type"`
	CategoryID uint     `json:"category_id"`
	Synonyms   []string `json:"synonyms"`
}

// CreateAtom creates a new atom
func (h *AtomHandler) CreateAtom(req CreateAtomRequest) AtomResponse {
	atom, err := h.service.CreateAtom(req.Value, req.Label, req.Type, req.CategoryID, req.Synonyms)
	if err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	return AtomResponse{Success: true, Data: atom}
}

// GetAtomByID retrieves an atom by ID
func (h *AtomHandler) GetAtomByID(id uint) AtomResponse {
	atom, err := h.service.GetAtomByID(id)
	if err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	return AtomResponse{Success: true, Data: atom}
}

// GetAtomsByCategoryRequest represents a request for atoms by category
type GetAtomsByCategoryRequest struct {
	CategoryID uint `json:"category_id"`
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
}

// GetAtomsByCategory retrieves atoms by category
func (h *AtomHandler) GetAtomsByCategory(req GetAtomsByCategoryRequest) AtomResponse {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 50
	}
	
	atoms, total, err := h.service.GetAtomsByCategory(req.CategoryID, req.Page, req.PageSize)
	if err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	
	return AtomResponse{Success: true, Data: map[string]interface{}{
		"atoms": atoms,
		"total": total,
		"page":  req.Page,
	}}
}

// UpdateAtomRequest represents an update atom request
type UpdateAtomRequest struct {
	ID       uint                   `json:"id"`
	Updates  map[string]interface{} `json:"updates"`
}

// UpdateAtom updates an atom
func (h *AtomHandler) UpdateAtom(req UpdateAtomRequest) AtomResponse {
	atom, err := h.service.UpdateAtom(req.ID, req.Updates)
	if err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	return AtomResponse{Success: true, Data: atom}
}

// DeleteAtom deletes an atom
func (h *AtomHandler) DeleteAtom(id uint) AtomResponse {
	if err := h.service.DeleteAtom(id); err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	return AtomResponse{Success: true}
}

// RecordUsage records atom usage
func (h *AtomHandler) RecordUsage(atomID uint) AtomResponse {
	if err := h.service.RecordUsage(atomID); err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	return AtomResponse{Success: true}
}

// FindAtomsBySynonym finds atoms by synonym search
func (h *AtomHandler) FindAtomsBySynonym(searchTerm string) AtomResponse {
	atoms, err := h.service.FindAtomsBySynonym(searchTerm)
	if err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	return AtomResponse{Success: true, Data: atoms}
}

// GetPopularAtoms gets popular atoms
func (h *AtomHandler) GetPopularAtoms(limit int) AtomResponse {
	if limit <= 0 {
		limit = 20
	}
	atoms, err := h.service.GetPopularAtoms(limit)
	if err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	return AtomResponse{Success: true, Data: atoms}
}

// BatchImportAtomsRequest represents a batch import request
type BatchImportAtomsRequest struct {
	JSONData string `json:"json_data"`
}

// BatchImportAtoms imports atoms from JSON
func (h *AtomHandler) BatchImportAtoms(req BatchImportAtomsRequest) AtomResponse {
	count, err := h.service.BatchImportAtoms(req.JSONData)
	if err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	return AtomResponse{Success: true, Data: map[string]int{"imported": count}}
}

// GetAllAtomsPaginated gets all atoms with pagination
func (h *AtomHandler) GetAllAtomsPaginated(page, pageSize int) AtomResponse {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 100
	}
	
	// Use search service to get all atoms
	atoms, err := h.service.FindAtomsBySynonym("")
	if err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	
	return AtomResponse{Success: true, Data: atoms}
}

// ExportAtoms exports all atoms as JSON
func (h *AtomHandler) ExportAtoms() AtomResponse {
	atoms, err := h.service.FindAtomsBySynonym("")
	if err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	
	jsonData, err := json.Marshal(atoms)
	if err != nil {
		return AtomResponse{Success: false, Error: err.Error()}
	}
	
	return AtomResponse{Success: true, Data: string(jsonData)}
}
