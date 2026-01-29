package handlers

import (
	"promptmaster/backend/utils"
)

// SeederHandler exposes seeder operations to the frontend
type SeederHandler struct {
	seeder *utils.Seeder
}

// NewSeederHandler creates a new SeederHandler
func NewSeederHandler(seeder *utils.Seeder) *SeederHandler {
	return &SeederHandler{seeder: seeder}
}

// SeederResponse represents the response structure
type SeederResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SeedAll seeds the database with default data
func (h *SeederHandler) SeedAll() SeederResponse {
	if err := h.seeder.SeedAll(); err != nil {
		return SeederResponse{Success: false, Error: err.Error()}
	}
	return SeederResponse{Success: true, Data: "Database seeded successfully"}
}

// ImportFromJSONRequest represents an import request
type ImportFromJSONRequest struct {
	JSONData string `json:"json_data"`
}

// ImportFromJSON imports data from JSON
func (h *SeederHandler) ImportFromJSON(req ImportFromJSONRequest) SeederResponse {
	if err := h.seeder.SeedFromJSON(req.JSONData); err != nil {
		return SeederResponse{Success: false, Error: err.Error()}
	}
	return SeederResponse{Success: true, Data: "Data imported successfully"}
}

// GetDefaultSeedData returns the default seed data as JSON
func (h *SeederHandler) GetDefaultSeedData() SeederResponse {
	data := h.seeder.GetDefaultSeedData()
	return SeederResponse{Success: true, Data: data}
}

// GetSeedStatus checks if database has been seeded
type GetSeedStatus struct {
	IsSeeded bool `json:"is_seeded"`
}

// GetSeedStatus returns the seed status of the database
func (h *SeederHandler) GetSeedStatus() SeederResponse {
	// Check if we have atoms
	data := h.seeder.GetDefaultSeedData()
	return SeederResponse{Success: true, Data: map[string]interface{}{
		"has_seed_data": len(data) > 0,
		"sample_data":   data[:min(200, len(data))],
	}}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
