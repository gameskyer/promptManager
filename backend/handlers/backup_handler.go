package handlers

import (
	"os"
	"promptmaster/backend/services"
)

// BackupHandler exposes backup operations to the frontend
type BackupHandler struct {
	service *services.BackupService
}

// NewBackupHandler creates a new BackupHandler
func NewBackupHandler(service *services.BackupService) *BackupHandler {
	return &BackupHandler{service: service}
}

// BackupResponse represents the response structure
type BackupResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ExportData exports all data to JSON
func (h *BackupHandler) ExportData() BackupResponse {
	jsonData, err := h.service.ExportToJSON()
	if err != nil {
		return BackupResponse{Success: false, Error: err.Error()}
	}
	return BackupResponse{Success: true, Data: jsonData}
}

// ImportDataRequest represents an import request
type ImportDataRequest struct {
	Data  string `json:"data"`  // JSON data to import
	Merge bool   `json:"merge"` // Whether to merge with existing data
}

// ImportData imports data from JSON
func (h *BackupHandler) ImportData(req ImportDataRequest) BackupResponse {
	if err := h.service.ImportJSON(req.Data, req.Merge); err != nil {
		return BackupResponse{Success: false, Error: err.Error()}
	}
	return BackupResponse{Success: true}
}

// ExportToFile exports data to a file and returns the path
func (h *BackupHandler) ExportToFile() BackupResponse {
	path, err := h.service.ExportAll()
	if err != nil {
		return BackupResponse{Success: false, Error: err.Error()}
	}
	return BackupResponse{Success: true, Data: path}
}

// GetBackupList returns list of available backups
func (h *BackupHandler) GetBackupList() BackupResponse {
	backups, err := h.service.GetBackupList()
	if err != nil {
		return BackupResponse{Success: false, Error: err.Error()}
	}
	return BackupResponse{Success: true, Data: backups}
}

// ExportToZip exports all data including images to a zip file
func (h *BackupHandler) ExportToZip() BackupResponse {
	path, err := h.service.ExportToZip()
	if err != nil {
		return BackupResponse{Success: false, Error: err.Error()}
	}
	return BackupResponse{Success: true, Data: path}
}

// ReadBackupFile reads a backup file and returns its content
func (h *BackupHandler) ReadBackupFile(filename string) BackupResponse {
	// Security: ensure filename doesn't contain path traversal
	if filename == "" || filename[0] == '/' || filename[0] == '\\' || filename == ".." {
		return BackupResponse{Success: false, Error: "invalid filename"}
	}
	
	backups, err := h.service.GetBackupList()
	if err != nil {
		return BackupResponse{Success: false, Error: err.Error()}
	}
	
	// Verify file is in backup list
	found := false
	for _, b := range backups {
		if b == filename {
			found = true
			break
		}
	}
	if !found {
		return BackupResponse{Success: false, Error: "backup not found"}
	}
	
	data, err := os.ReadFile(filename)
	if err != nil {
		return BackupResponse{Success: false, Error: err.Error()}
	}
	
	return BackupResponse{Success: true, Data: string(data)}
}
