package services

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"promptmaster/backend/config"
	"promptmaster/backend/models"
	"time"

	"gorm.io/gorm"
)

// BackupService handles data export and import
type BackupService struct {
	db *gorm.DB
}

// NewBackupService creates a new BackupService
func NewBackupService(db *gorm.DB) *BackupService {
	return &BackupService{db: db}
}

// AtomExport represents atom for export (without relationships)
type AtomExport struct {
	ID          uint              `json:"id"`
	Value       string            `json:"value"`
	Label       string            `json:"label"`
	Synonyms    models.StringSlice `json:"synonyms"`
	Type        string            `json:"type"`
	CategoryID  uint              `json:"category_id"`
	UsageCount  int               `json:"usage_count"`
	LastUsedAt  *time.Time        `json:"last_used_at"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// PresetVersionExport represents version for export
type PresetVersionExport struct {
	ID            uint      `json:"id"`
	PresetID      uint      `json:"preset_id"`
	VersionNum    int       `json:"version_num"`
	Snapshot      models.JSON `json:"snapshot"`
	ThumbnailPath string    `json:"thumbnail_path"`
	IsStarred     bool      `json:"is_starred"`
	CreatedAt     time.Time `json:"created_at"`
	DiffStats     string    `json:"diff_stats"`
}

// PresetExport represents preset for export
type PresetExport struct {
	ID             int                   `json:"id"`
	Title          string                `json:"title"`
	CurrentVersion int                   `json:"current_version"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	IsDeleted      bool                  `json:"is_deleted"`
	Versions       []PresetVersionExport `json:"versions"`
}

// ExportData represents all exportable data
type ExportData struct {
	Version    string            `json:"version"`
	ExportTime time.Time         `json:"export_time"`
	Categories []models.Category `json:"categories"`
	Atoms      []AtomExport      `json:"atoms"`
	Presets    []PresetExport    `json:"presets"`
}

// ExportAll exports all data to a JSON file
func (s *BackupService) ExportAll() (string, error) {
	data, err := s.exportData()
	if err != nil {
		return "", err
	}

	// Create export file
	exportFileName := fmt.Sprintf("promptmaster_backup_%s.json", time.Now().Format("20060102_150405"))
	exportPath := filepath.Join(config.BackupDir, exportFileName)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(exportPath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("failed to write export file: %w", err)
	}

	return exportPath, nil
}

// ExportToJSON exports data to a JSON string
func (s *BackupService) ExportToJSON() (string, error) {
	data, err := s.exportData()
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	return string(jsonData), nil
}

// exportData exports all data
func (s *BackupService) exportData() (*ExportData, error) {
	var data ExportData
	data.Version = config.AppVersion
	data.ExportTime = time.Now()

	// Export categories
	if err := s.db.Find(&data.Categories).Error; err != nil {
		return nil, fmt.Errorf("failed to export categories: %w", err)
	}

	// Export atoms (convert to export format)
	var atoms []models.Atom
	if err := s.db.Find(&atoms).Error; err != nil {
		return nil, fmt.Errorf("failed to export atoms: %w", err)
	}
	for _, atom := range atoms {
		data.Atoms = append(data.Atoms, AtomExport{
			ID:         atom.ID,
			Value:      atom.Value,
			Label:      atom.Label,
			Synonyms:   atom.Synonyms,
			Type:       atom.Type,
			CategoryID: atom.CategoryID,
			UsageCount: atom.UsageCount,
			LastUsedAt: atom.LastUsedAt,
			CreatedAt:  atom.CreatedAt,
			UpdatedAt:  atom.UpdatedAt,
		})
	}

	// Export presets with versions
	var presets []models.Preset
	if err := s.db.Preload("Versions").Find(&presets).Error; err != nil {
		return nil, fmt.Errorf("failed to export presets: %w", err)
	}

	for _, preset := range presets {
		presetExport := PresetExport{
			ID:             int(preset.ID),
			Title:          preset.Title,
			CurrentVersion: preset.CurrentVersion,
			CreatedAt:      preset.CreatedAt,
			UpdatedAt:      preset.UpdatedAt,
			IsDeleted:      preset.IsDeleted,
		}
		for _, v := range preset.Versions {
			presetExport.Versions = append(presetExport.Versions, PresetVersionExport{
				ID:            v.ID,
				PresetID:      v.PresetID,
				VersionNum:    v.VersionNum,
				Snapshot:      v.Snapshot,
				ThumbnailPath: v.ThumbnailPath,
				IsStarred:     v.IsStarred,
				CreatedAt:     v.CreatedAt,
				DiffStats:     v.DiffStats,
			})
		}
		data.Presets = append(data.Presets, presetExport)
	}

	return &data, nil
}

// ImportAll imports data from JSON
func (s *BackupService) ImportAll(jsonData string) error {
	var data ExportData
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return fmt.Errorf("failed to parse import data: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Import categories
		for _, cat := range data.Categories {
			var existing models.Category
			if err := tx.First(&existing, cat.ID).Error; err != nil {
				// Create new
				if err := tx.Create(&cat).Error; err != nil {
					return fmt.Errorf("failed to import category %d: %w", cat.ID, err)
				}
			} else {
				// Update existing
				if err := tx.Model(&existing).Updates(cat).Error; err != nil {
					return fmt.Errorf("failed to update category %d: %w", cat.ID, err)
				}
			}
		}

		// Import atoms
		for _, atomExport := range data.Atoms {
			atom := models.Atom{
				ID:          atomExport.ID,
				Value:       atomExport.Value,
				Label:       atomExport.Label,
				Synonyms:    atomExport.Synonyms,
				Type:        atomExport.Type,
				CategoryID:  atomExport.CategoryID,
				UsageCount:  atomExport.UsageCount,
				LastUsedAt:  atomExport.LastUsedAt,
				CreatedAt:   atomExport.CreatedAt,
				UpdatedAt:   atomExport.UpdatedAt,
			}
			var existing models.Atom
			if err := tx.First(&existing, atom.ID).Error; err != nil {
				if err := tx.Create(&atom).Error; err != nil {
					return fmt.Errorf("failed to import atom %d: %w", atom.ID, err)
				}
			} else {
				if err := tx.Model(&existing).Updates(map[string]interface{}{
					"value":       atom.Value,
					"label":       atom.Label,
					"synonyms":    atom.Synonyms,
					"type":        atom.Type,
					"category_id": atom.CategoryID,
					"usage_count": atom.UsageCount,
					"updated_at":  atom.UpdatedAt,
				}).Error; err != nil {
					return fmt.Errorf("failed to update atom %d: %w", atom.ID, err)
				}
			}
		}

		// Import presets with versions
		for _, presetExport := range data.Presets {
			preset := models.Preset{
				ID:             uint(presetExport.ID),
				Title:          presetExport.Title,
				CurrentVersion: presetExport.CurrentVersion,
				CreatedAt:      presetExport.CreatedAt,
				UpdatedAt:      presetExport.UpdatedAt,
				IsDeleted:      presetExport.IsDeleted,
			}

			var existing models.Preset
			if err := tx.First(&existing, preset.ID).Error; err != nil {
				// Create new preset
				if err := tx.Create(&preset).Error; err != nil {
					return fmt.Errorf("failed to import preset %d: %w", preset.ID, err)
				}
			} else {
				// Update existing
				if err := tx.Model(&existing).Updates(map[string]interface{}{
					"title":           preset.Title,
					"current_version": preset.CurrentVersion,
					"is_deleted":      preset.IsDeleted,
					"updated_at":      preset.UpdatedAt,
				}).Error; err != nil {
					return fmt.Errorf("failed to update preset %d: %w", preset.ID, err)
				}
			}

			// Import versions
			for _, v := range presetExport.Versions {
				version := models.PresetVersion{
					ID:            v.ID,
					PresetID:      v.PresetID,
					VersionNum:    v.VersionNum,
					Snapshot:      v.Snapshot,
					ThumbnailPath: v.ThumbnailPath,
					IsStarred:     v.IsStarred,
					CreatedAt:     v.CreatedAt,
					DiffStats:     v.DiffStats,
				}
				var existingVersion models.PresetVersion
				if err := tx.First(&existingVersion, version.ID).Error; err != nil {
					if err := tx.Create(&version).Error; err != nil {
						return fmt.Errorf("failed to import version %d: %w", version.ID, err)
					}
				} else {
					if err := tx.Model(&existingVersion).Updates(version).Error; err != nil {
						return fmt.Errorf("failed to update version %d: %w", version.ID, err)
					}
				}
			}
		}

		return nil
	})
}

// ImportJSON imports data from JSON string
func (s *BackupService) ImportJSON(jsonStr string, merge bool) error {
	var data ExportData
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	if !merge {
		// Clear existing data if not merging
		s.db.Exec("DELETE FROM preset_versions")
		s.db.Exec("DELETE FROM presets")
		s.db.Exec("DELETE FROM atoms")
		s.db.Exec("DELETE FROM categories")
	}

	return s.ImportAll(jsonStr)
}

// GetBackupList returns list of available backups
func (s *BackupService) GetBackupList() ([]string, error) {
	entries, err := os.ReadDir(config.BackupDir)
	if err != nil {
		return nil, err
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}

// ExportToZip exports all data including images to a zip file
func (s *BackupService) ExportToZip() (string, error) {
	// Create export JSON
	jsonData, err := s.ExportToJSON()
	if err != nil {
		return "", err
	}

	// Create zip file
	zipFileName := fmt.Sprintf("promptmaster_backup_%s.zip", time.Now().Format("20060102_150405"))
	zipPath := filepath.Join(config.BackupDir, zipFileName)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add data.json to zip
	dataWriter, err := zipWriter.Create("data.json")
	if err != nil {
		return "", fmt.Errorf("failed to create data entry: %w", err)
	}
	if _, err := dataWriter.Write([]byte(jsonData)); err != nil {
		return "", fmt.Errorf("failed to write data: %w", err)
	}

	// Add images
	imageEntries, err := os.ReadDir(config.ImageDir)
	if err == nil {
		for _, entry := range imageEntries {
			if entry.IsDir() {
				continue
			}
			filePath := filepath.Join(config.ImageDir, entry.Name())
			file, err := os.Open(filePath)
			if err != nil {
				continue
			}
			defer file.Close()

			writer, err := zipWriter.Create(filepath.Join("images", entry.Name()))
			if err != nil {
				continue
			}
			io.Copy(writer, file)
		}
	}

	return zipPath, nil
}
