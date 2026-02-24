package services

import (
	"promptmaster/backend/models"
	"time"

	"gorm.io/gorm"
)

// BatchService handles batch operations
type BatchService struct {
	db *gorm.DB
}

// NewBatchService creates a new BatchService
func NewBatchService(db *gorm.DB) *BatchService {
	return &BatchService{db: db}
}

// BatchMoveCategory moves multiple atoms to a new category
func (s *BatchService) BatchMoveCategory(atomIDs []uint, categoryID uint) (int, error) {
	if len(atomIDs) == 0 {
		return 0, nil
	}

	result := s.db.Model(&models.Atom{}).
		Where("id IN ?", atomIDs).
		Updates(map[string]interface{}{
			"category_id": categoryID,
			"updated_at":  time.Now(),
		})

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// BatchUpdateType updates the type of multiple atoms
func (s *BatchService) BatchUpdateType(atomIDs []uint, atomType string) (int, error) {
	if len(atomIDs) == 0 {
		return 0, nil
	}

	result := s.db.Model(&models.Atom{}).
		Where("id IN ?", atomIDs).
		Updates(map[string]interface{}{
			"type":       atomType,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// BatchDelete soft deletes multiple atoms
func (s *BatchService) BatchDelete(atomIDs []uint) (int, error) {
	if len(atomIDs) == 0 {
		return 0, nil
	}

	result := s.db.Delete(&models.Atom{}, "id IN ?", atomIDs)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// BatchAddSynonyms adds synonyms to multiple atoms
func (s *BatchService) BatchAddSynonyms(atomIDs []uint, synonyms []string) (int, error) {
	if len(atomIDs) == 0 || len(synonyms) == 0 {
		return 0, nil
	}

	var atoms []models.Atom
	if err := s.db.Where("id IN ?", atomIDs).Find(&atoms).Error; err != nil {
		return 0, err
	}

	updatedCount := 0
	for _, atom := range atoms {
		// Merge existing synonyms with new ones
		synonymMap := make(map[string]bool)
		for _, s := range atom.Synonyms {
			synonymMap[s] = true
		}
		for _, s := range synonyms {
			synonymMap[s] = true
		}

		var mergedSynonyms []string
		for s := range synonymMap {
			mergedSynonyms = append(mergedSynonyms, s)
		}

		result := s.db.Model(&atom).Updates(map[string]interface{}{
			"synonyms":   models.StringSlice(mergedSynonyms),
			"updated_at": time.Now(),
		})

		if result.Error == nil {
			updatedCount++
		}
	}

	return updatedCount, nil
}

// BatchClearCategory clears the category of multiple atoms
func (s *BatchService) BatchClearCategory(atomIDs []uint) (int, error) {
	if len(atomIDs) == 0 {
		return 0, nil
	}

	result := s.db.Model(&models.Atom{}).
		Where("id IN ?", atomIDs).
		Updates(map[string]interface{}{
			"category_id": 0,
			"updated_at":  time.Now(),
		})

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// BatchOperationResult represents the result of a batch operation
type BatchOperationResult struct {
	SuccessCount int      `json:"success_count"`
	FailedIDs    []uint   `json:"failed_ids,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}
