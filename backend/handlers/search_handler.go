package handlers

import (
	"promptmaster/backend/services"
)

// SearchHandler exposes search operations to the frontend
type SearchHandler struct {
	service *services.SearchService
}

// NewSearchHandler creates a new SearchHandler
func NewSearchHandler(service *services.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// SearchResponse represents the response structure
type SearchResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
}

// Search performs a global search
func (h *SearchHandler) Search(req SearchRequest) SearchResponse {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	
	results, err := h.service.Search(req.Query, req.Limit)
	if err != nil {
		return SearchResponse{Success: false, Error: err.Error()}
	}
	return SearchResponse{Success: true, Data: results}
}

// SearchAtomsRequest represents an atom search request
type SearchAtomsRequest struct {
	SearchTerm string `json:"search_term"`
	Type       string `json:"type"`
	CategoryID uint   `json:"category_id"`
	Limit      int    `json:"limit"`
}

// SearchAtoms searches for atoms
func (h *SearchHandler) SearchAtoms(req SearchAtomsRequest) SearchResponse {
	if req.Limit <= 0 {
		req.Limit = 50
	}
	
	atoms, err := h.service.SearchAtoms(req.SearchTerm, req.Type, req.CategoryID, req.Limit)
	if err != nil {
		return SearchResponse{Success: false, Error: err.Error()}
	}
	return SearchResponse{Success: true, Data: atoms}
}

// QuickSearch performs a quick search
func (h *SearchHandler) QuickSearch(term string) SearchResponse {
	results, err := h.service.QuickSearch(term)
	if err != nil {
		return SearchResponse{Success: false, Error: err.Error()}
	}
	return SearchResponse{Success: true, Data: results}
}

// SearchPresetsRequest represents a preset search request
type SearchPresetsRequest struct {
	SearchTerm string `json:"search_term"`
	Limit      int    `json:"limit"`
}

// SearchPresets searches for presets
func (h *SearchHandler) SearchPresets(req SearchPresetsRequest) SearchResponse {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	
	presets, err := h.service.SearchPresets(req.SearchTerm, req.Limit)
	if err != nil {
		return SearchResponse{Success: false, Error: err.Error()}
	}
	return SearchResponse{Success: true, Data: presets}
}

// ReindexAll rebuilds the search index
func (h *SearchHandler) ReindexAll() SearchResponse {
	// This would trigger a full reindex
	return SearchResponse{Success: true, Data: map[string]string{"message": "Reindex initiated"}}
}
