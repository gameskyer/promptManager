package services

import (
	"encoding/json"
	"fmt"
	"promptmaster/backend/models"
	"time"

	"gorm.io/gorm"
)

// VersionService handles version control operations
type VersionService struct {
	db *gorm.DB
}

// NewVersionService creates a new VersionService
func NewVersionService(db *gorm.DB) *VersionService {
	return &VersionService{db: db}
}

// VersionInput represents the input for creating a new version
type VersionInput struct {
	PresetID      uint
	PosText       string
	NegText       string
	AtomIDs       []uint
	Params        map[string]interface{}
	PreviewPaths  []string
	ThumbnailPath string
}

// CreateVersion creates a new version for a preset
func (s *VersionService) CreateVersion(input VersionInput) (*models.PresetVersion, error) {
	var newVersionNum int
	
	// Get current max version number
	var maxVersion models.PresetVersion
	err := s.db.Where("preset_id = ?", input.PresetID).
		Order("version_num DESC").
		First(&maxVersion).Error
	
	if err == gorm.ErrRecordNotFound {
		newVersionNum = 1
	} else if err != nil {
		return nil, err
	} else {
		newVersionNum = maxVersion.VersionNum + 1
	}
	
	// Build snapshot
	snapshot := &models.SnapshotData{
		PosText:      input.PosText,
		NegText:      input.NegText,
		AtomIDs:      input.AtomIDs,
		Params:       input.Params,
		PreviewPaths: input.PreviewPaths,
	}
	
	// Calculate diff if there's a previous version
	diffStats := ""
	if newVersionNum > 1 {
		var prevVersion models.PresetVersion
		if err := s.db.Where("preset_id = ? AND version_num = ?", 
			input.PresetID, newVersionNum-1).
			First(&prevVersion).Error; err == nil {
			diffStats = s.calculateDiff(&prevVersion, snapshot)
		}
	} else {
		diffStats = "+0/-0"
	}
	
	version := &models.PresetVersion{
		PresetID:      input.PresetID,
		VersionNum:    newVersionNum,
		ThumbnailPath: input.ThumbnailPath,
		IsStarred:     false,
		CreatedAt:     time.Now(),
		DiffStats:     diffStats,
	}
	version.SetSnapshotData(snapshot)
	
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(version).Error; err != nil {
			return err
		}
		
		// Update preset's current version
		return tx.Model(&models.Preset{}).
			Where("id = ?", input.PresetID).
			Updates(map[string]interface{}{
				"current_version": newVersionNum,
				"updated_at":      time.Now(),
			}).Error
	})
	
	if err != nil {
		return nil, err
	}
	
	return version, nil
}

// GetVersion retrieves a specific version
func (s *VersionService) GetVersion(presetID uint, versionNum int) (*models.PresetVersion, error) {
	var version models.PresetVersion
	if err := s.db.Where("preset_id = ? AND version_num = ?", presetID, versionNum).
		First(&version).Error; err != nil {
		return nil, err
	}
	return &version, nil
}

// GetVersionHistory retrieves all versions for a preset
func (s *VersionService) GetVersionHistory(presetID uint, limit int) ([]models.PresetVersion, error) {
	var versions []models.PresetVersion
	query := s.db.Where("preset_id = ?", presetID).Order("version_num DESC").Preload("Previews")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

// GetLatestVersions retrieves the latest N versions
func (s *VersionService) GetLatestVersions(presetID uint, count int) ([]models.PresetVersion, error) {
	var versions []models.PresetVersion
	if err := s.db.Where("preset_id = ?", presetID).
		Order("version_num DESC").
		Limit(count).
		Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

// StarVersion toggles the star status of a version
func (s *VersionService) StarVersion(versionID uint, starred bool) error {
	return s.db.Model(&models.PresetVersion{}).
		Where("id = ?", versionID).
		Update("is_starred", starred).Error
}

// RollbackToVersion creates a new version based on an old version
func (s *VersionService) RollbackToVersion(presetID uint, targetVersionNum int) (*models.PresetVersion, error) {
	var targetVersion models.PresetVersion
	if err := s.db.Where("preset_id = ? AND version_num = ?", presetID, targetVersionNum).
		First(&targetVersion).Error; err != nil {
		return nil, fmt.Errorf("target version not found: %w", err)
	}
	
	snapshot, err := targetVersion.ToSnapshotData()
	if err != nil {
		return nil, err
	}
	
	input := VersionInput{
		PresetID:      presetID,
		PosText:       snapshot.PosText,
		NegText:       snapshot.NegText,
		AtomIDs:       snapshot.AtomIDs,
		Params:        snapshot.Params,
		PreviewPaths:  snapshot.PreviewPaths,
		ThumbnailPath: targetVersion.ThumbnailPath,
	}
	
	return s.CreateVersion(input)
}

// CompareVersions compares two versions and returns detailed diff
func (s *VersionService) CompareVersions(presetID uint, v1, v2 int) (*VersionDiff, error) {
	version1, err := s.GetVersion(presetID, v1)
	if err != nil {
		return nil, fmt.Errorf("version %d not found: %w", v1, err)
	}
	
	version2, err := s.GetVersion(presetID, v2)
	if err != nil {
		return nil, fmt.Errorf("version %d not found: %w", v2, err)
	}
	
	data1, err := version1.ToSnapshotData()
	if err != nil {
		return nil, err
	}
	
	data2, err := version2.ToSnapshotData()
	if err != nil {
		return nil, err
	}
	
	return s.buildDetailedDiff(version1, version2, data1, data2)
}

// VersionDiff represents a detailed version comparison
type VersionDiff struct {
	Version1     *models.PresetVersion `json:"version1"`
	Version2     *models.PresetVersion `json:"version2"`
	AddedAtoms   []models.Atom         `json:"added_atoms"`
	RemovedAtoms []models.Atom         `json:"removed_atoms"`
	ParamChanges []ParamChange         `json:"param_changes"`
	PosTextDiff  TextDiff              `json:"pos_text_diff"`
	NegTextDiff  TextDiff              `json:"neg_text_diff"`
}

// ParamChange represents a parameter change
type ParamChange struct {
	Param    string      `json:"param"`
	OldValue interface{} `json:"old_value"`
	NewValue interface{} `json:"new_value"`
}

// TextDiff represents text differences
type TextDiff struct {
	OldText string `json:"old_text"`
	NewText string `json:"new_text"`
}

func (s *VersionService) buildDetailedDiff(v1, v2 *models.PresetVersion, 
	d1, d2 *models.SnapshotData) (*VersionDiff, error) {
	
	diff := &VersionDiff{
		Version1: v1,
		Version2: v2,
	}
	
	// Find added and removed atoms
	oldAtomMap := make(map[uint]bool)
	for _, id := range d1.AtomIDs {
		oldAtomMap[id] = true
	}
	
	newAtomMap := make(map[uint]bool)
	for _, id := range d2.AtomIDs {
		newAtomMap[id] = true
	}
	
	// Get added atoms
	var addedIDs []uint
	for id := range newAtomMap {
		if !oldAtomMap[id] {
			addedIDs = append(addedIDs, id)
		}
	}
	
	// Get removed atoms
	var removedIDs []uint
	for id := range oldAtomMap {
		if !newAtomMap[id] {
			removedIDs = append(removedIDs, id)
		}
	}
	
	// Fetch atom details
	if len(addedIDs) > 0 {
		s.db.Where("id IN ?", addedIDs).Find(&diff.AddedAtoms)
	}
	if len(removedIDs) > 0 {
		s.db.Where("id IN ?", removedIDs).Find(&diff.RemovedAtoms)
	}
	
	// Find parameter changes
	for key, newVal := range d2.Params {
		if oldVal, exists := d1.Params[key]; exists {
			if !isEqual(oldVal, newVal) {
				diff.ParamChanges = append(diff.ParamChanges, ParamChange{
					Param:    key,
					OldValue: oldVal,
					NewValue: newVal,
				})
			}
		} else {
			diff.ParamChanges = append(diff.ParamChanges, ParamChange{
				Param:    key,
				OldValue: nil,
				NewValue: newVal,
			})
		}
	}
	
	for key, oldVal := range d1.Params {
		if _, exists := d2.Params[key]; !exists {
			diff.ParamChanges = append(diff.ParamChanges, ParamChange{
				Param:    key,
				OldValue: oldVal,
				NewValue: nil,
			})
		}
	}
	
	diff.PosTextDiff = TextDiff{OldText: d1.PosText, NewText: d2.PosText}
	diff.NegTextDiff = TextDiff{OldText: d1.NegText, NewText: d2.NegText}
	
	return diff, nil
}

func (s *VersionService) calculateDiff(oldVersion *models.PresetVersion, newSnapshot *models.SnapshotData) string {
	oldData, err := oldVersion.ToSnapshotData()
	if err != nil {
		return "+0/-0"
	}
	
	// Compare atom IDs
	oldAtomMap := make(map[uint]bool)
	for _, id := range oldData.AtomIDs {
		oldAtomMap[id] = true
	}
	
	newAtomMap := make(map[uint]bool)
	for _, id := range newSnapshot.AtomIDs {
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
	
	return fmt.Sprintf("+%d/-%d", added, removed)
}

func isEqual(a, b interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}

// DeleteVersion deletes a specific version
func (s *VersionService) DeleteVersion(versionID uint) error {
	return s.db.Delete(&models.PresetVersion{}, versionID).Error
}

// GetStarredVersions returns all starred versions for a preset
func (s *VersionService) GetStarredVersions(presetID uint) ([]models.PresetVersion, error) {
	var versions []models.PresetVersion
	if err := s.db.Where("preset_id = ? AND is_starred = ?", presetID, true).
		Order("version_num DESC").
		Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}
