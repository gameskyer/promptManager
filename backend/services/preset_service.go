package services

import (
	"fmt"
	"promptmaster/backend/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// PresetService handles preset-related operations
type PresetService struct {
	db          *gorm.DB
	imageService *ImageService
}

// NewPresetService creates a new PresetService
func NewPresetService(db *gorm.DB) *PresetService {
	return &PresetService{
		db:          db,
		imageService: NewImageService(db),
	}
}

// CreatePreset creates a new preset with initial version
func (s *PresetService) CreatePreset(title string, categoryID uint, posText, negText string, atomIDs []uint, params map[string]interface{}, previewBase64s []string) (*models.Preset, error) {
	preset := &models.Preset{
		Title:      title,
		CategoryID: categoryID,
	}
	
	var previewPaths []string
	
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(preset).Error; err != nil {
			return err
		}
		
		// Upload preview images
		if len(previewBase64s) > 0 {
			for _, base64Data := range previewBase64s {
				result, err := s.imageService.UploadImage(UploadImageRequest{
					Data:     base64Data,
					PresetID: preset.ID,
				})
				if err != nil {
					// Log error but continue
					fmt.Printf("Failed to upload image: %v\n", err)
					continue
				}
				previewPaths = append(previewPaths, result.FilePath)
			}
		}
		
		// Create initial version (V1)
		snapshot := &models.SnapshotData{
			PosText:      posText,
			NegText:      negText,
			Params:       params,
			AtomIDs:      atomIDs,
			PreviewPaths: previewPaths,
		}
		
		version := &models.PresetVersion{
			PresetID:   preset.ID,
			VersionNum: 1,
			CreatedAt:  time.Now(),
		}
		version.SetSnapshotData(snapshot)
		
		if err := tx.Create(version).Error; err != nil {
			return err
		}
		
		// Update preset's current version
		preset.CurrentVersion = 1
		return tx.Model(preset).Update("current_version", 1).Error
	})
	
	if err != nil {
		return nil, err
	}
	
	return preset, nil
}

// GetPresetByID retrieves a preset by ID
func (s *PresetService) GetPresetByID(id uint) (*models.Preset, error) {
	var preset models.Preset
	if err := s.db.Preload("Versions").First(&preset, id).Error; err != nil {
		return nil, err
	}
	return &preset, nil
}

// GetPresets retrieves presets with pagination
func (s *PresetService) GetPresets(page, pageSize int, categoryID uint, includeDeleted bool) ([]models.Preset, int64, error) {
	var presets []models.Preset
	var total int64
	
	query := s.db.Model(&models.Preset{})
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	
	// Filter by category if specified
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	
	query.Count(&total)
	
	offset := (page - 1) * pageSize
	if err := query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&presets).Error; err != nil {
		return nil, 0, err
	}
	
	return presets, total, nil
}

// GetPresetsByCategory retrieves presets by category ID
func (s *PresetService) GetPresetsByCategory(categoryID uint, page, pageSize int) ([]models.Preset, int64, error) {
	return s.GetPresets(page, pageSize, categoryID, false)
}

// UpdatePreset updates a preset's basic info
func (s *PresetService) UpdatePreset(id uint, title string, categoryID uint) (*models.Preset, error) {
	var preset models.Preset
	if err := s.db.First(&preset, id).Error; err != nil {
		return nil, err
	}
	
	updates := map[string]interface{}{
		"title":       title,
		"category_id": categoryID,
		"updated_at":  time.Now(),
	}
	
	if err := s.db.Model(&preset).Updates(updates).Error; err != nil {
		return nil, err
	}
	
	return &preset, nil
}

// SoftDeletePreset soft deletes a preset
func (s *PresetService) SoftDeletePreset(id uint) error {
	return s.db.Model(&models.Preset{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}).Error
}

// RestorePreset restores a soft-deleted preset
func (s *PresetService) RestorePreset(id uint) error {
	return s.db.Model(&models.Preset{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": false,
		"updated_at": time.Now(),
	}).Error
}

// BuildPromptText builds the final prompt text from atoms
func (s *PresetService) BuildPromptText(atomIDs []uint) (string, []uint, error) {
	if len(atomIDs) == 0 {
		return "", []uint{}, nil
	}
	
	var atoms []models.Atom
	if err := s.db.Where("id IN ?", atomIDs).Find(&atoms).Error; err != nil {
		return "", nil, err
	}
	
	// Build text in order of atomIDs
	atomMap := make(map[uint]models.Atom)
	for _, atom := range atoms {
		atomMap[atom.ID] = atom
	}
	
	var parts []string
	var foundIDs []uint
	for _, id := range atomIDs {
		if atom, ok := atomMap[id]; ok {
			parts = append(parts, atom.Value)
			foundIDs = append(foundIDs, id)
		}
	}
	
	return strings.Join(parts, ", "), foundIDs, nil
}

// GetCurrentWorkState retrieves the current working state for a preset
func (s *PresetService) GetCurrentWorkState(presetID uint) (*models.SnapshotData, int, error) {
	var version models.PresetVersion
	if err := s.db.Where("preset_id = ? AND version_num = ?", 
		presetID, s.db.Model(&models.Preset{}).Select("current_version").Where("id = ?", presetID)).
		First(&version).Error; err != nil {
		return nil, 0, err
	}
	
	data, err := version.ToSnapshotData()
	if err != nil {
		return nil, 0, err
	}
	
	return data, version.VersionNum, nil
}

// ForkPreset creates a new preset based on an existing version
func (s *PresetService) ForkPreset(presetID uint, versionNum int, newTitle string, categoryID uint) (*models.Preset, error) {
	var sourceVersion models.PresetVersion
	if err := s.db.Where("preset_id = ? AND version_num = ?", presetID, versionNum).
		First(&sourceVersion).Error; err != nil {
		return nil, err
	}
	
	snapshot, err := sourceVersion.ToSnapshotData()
	if err != nil {
		return nil, err
	}
	
	return s.CreatePreset(newTitle, categoryID, snapshot.PosText, snapshot.NegText, snapshot.AtomIDs, snapshot.Params, snapshot.PreviewPaths)
}

// CleanupOldVersions removes old versions based on retention policy
func (s *PresetService) CleanupOldVersions(presetID uint, keepCount int) error {
	var versions []models.PresetVersion
	s.db.Where("preset_id = ? AND is_starred = ?", presetID, false).
		Order("version_num DESC").
		Offset(keepCount).
		Find(&versions)
	
	for _, v := range versions {
		s.db.Delete(&v)
	}
	
	return nil
}

// GenerateVersionDiff calculates diff between two versions
func (s *PresetService) GenerateVersionDiff(oldVersion, newVersion *models.PresetVersion) (string, error) {
	oldData, err := oldVersion.ToSnapshotData()
	if err != nil {
		return "", err
	}
	
	newData, err := newVersion.ToSnapshotData()
	if err != nil {
		return "", err
	}
	
	// Compare atom IDs
	oldAtomMap := make(map[uint]bool)
	for _, id := range oldData.AtomIDs {
		oldAtomMap[id] = true
	}
	
	newAtomMap := make(map[uint]bool)
	for _, id := range newData.AtomIDs {
		newAtomMap[id] = true
	}
	
	added := 0
	for id := range newAtomMap {
		if !oldAtomMap[id] {
			added++
		}
	}
	
	removed := 0
	for id := range oldAtomMap {
		if !newAtomMap[id] {
			removed++
		}
	}
	
	if added == 0 && removed == 0 {
		return "+0/-0", nil
	}
	
	return fmt.Sprintf("+%d/-%d", added, removed), nil
}

// PresetWithSnapshot is a preset with its current version snapshot data
type PresetWithSnapshot struct {
	models.Preset
	PosText   string                 `json:"pos_text"`
	NegText   string                 `json:"neg_text"`
	Params    map[string]interface{} `json:"params"`
	AtomIDs   []uint                 `json:"atom_ids"`
	Previews  []string               `json:"previews"`
	Thumbnail string                 `json:"thumbnail"`
	Loras     []map[string]interface{} `json:"loras"`
}

// ToPresetWithSnapshot converts a preset to include current version snapshot data
func (s *PresetService) ToPresetWithSnapshot(preset *models.Preset) (*PresetWithSnapshot, error) {
	result := &PresetWithSnapshot{
		Preset:   *preset,
		Params:   make(map[string]interface{}),
		AtomIDs:  []uint{},
		Previews: []string{},
		Loras:    []map[string]interface{}{},
	}
	
	// Find current version
	var currentVersion *models.PresetVersion
	for _, v := range preset.Versions {
		if v.VersionNum == preset.CurrentVersion {
			currentVersion = &v
			break
		}
	}
	
	if currentVersion == nil && len(preset.Versions) > 0 {
		currentVersion = &preset.Versions[len(preset.Versions)-1]
	}
	
	if currentVersion != nil {
		snapshot, err := currentVersion.ToSnapshotData()
		if err == nil {
			result.PosText = snapshot.PosText
			result.NegText = snapshot.NegText
			result.Params = snapshot.Params
			result.AtomIDs = snapshot.AtomIDs
			result.Previews = snapshot.PreviewPaths
			if len(snapshot.PreviewPaths) > 0 {
				result.Thumbnail = snapshot.PreviewPaths[0]
			}
			// Extract LoRAs from params
			if loraData, ok := snapshot.Params["loras"].([]interface{}); ok {
				for _, l := range loraData {
					if lora, ok := l.(map[string]interface{}); ok {
						result.Loras = append(result.Loras, lora)
					}
				}
			}
		}
	}
	
	return result, nil
}
