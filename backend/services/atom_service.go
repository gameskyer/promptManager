package services

import (
	"encoding/json"
	"promptmaster/backend/models"
	"time"

	"gorm.io/gorm"
)

// AtomService handles atom-related operations
type AtomService struct {
	db *gorm.DB
}

// NewAtomService creates a new AtomService
func NewAtomService(db *gorm.DB) *AtomService {
	return &AtomService{db: db}
}

// CreateAtom creates a new atom
func (s *AtomService) CreateAtom(value, label, atomType string, categoryID uint, synonyms []string) (*models.Atom, error) {
	atom := &models.Atom{
		Value:      value,
		Label:      label,
		Type:       atomType,
		CategoryID: categoryID,
		Synonyms:   models.StringSlice(synonyms),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	if err := s.db.Create(atom).Error; err != nil {
		return nil, err
	}
	
	return atom, nil
}

// GetAtomByID retrieves an atom by ID
func (s *AtomService) GetAtomByID(id uint) (*models.Atom, error) {
	var atom models.Atom
	if err := s.db.Preload("Category").First(&atom, id).Error; err != nil {
		return nil, err
	}
	return &atom, nil
}

// GetAtomByValue retrieves an atom by its English value
func (s *AtomService) GetAtomByValue(value string) (*models.Atom, error) {
	var atom models.Atom
	if err := s.db.Where("value = ?", value).First(&atom).Error; err != nil {
		return nil, err
	}
	return &atom, nil
}

// GetAtomsByCategory retrieves atoms by category ID
func (s *AtomService) GetAtomsByCategory(categoryID uint, page, pageSize int) ([]models.Atom, int64, error) {
	var atoms []models.Atom
	var total int64
	
	query := s.db.Model(&models.Atom{}).Where("category_id = ?", categoryID)
	query.Count(&total)
	
	offset := (page - 1) * pageSize
	if err := query.Order("usage_count DESC").Offset(offset).Limit(pageSize).Find(&atoms).Error; err != nil {
		return nil, 0, err
	}
	
	return atoms, total, nil
}

// UpdateAtom updates an atom
func (s *AtomService) UpdateAtom(id uint, updates map[string]interface{}) (*models.Atom, error) {
	var atom models.Atom
	if err := s.db.First(&atom, id).Error; err != nil {
		return nil, err
	}
	
	if err := s.db.Model(&atom).Updates(updates).Error; err != nil {
		return nil, err
	}
	
	return &atom, nil
}

// DeleteAtom soft deletes an atom
func (s *AtomService) DeleteAtom(id uint) error {
	return s.db.Delete(&models.Atom{}, id).Error
}

// RecordUsage records that an atom was used
func (s *AtomService) RecordUsage(atomID uint) error {
	now := time.Now()
	return s.db.Model(&models.Atom{}).Where("id = ?", atomID).Updates(map[string]interface{}{
		"usage_count":  gorm.Expr("usage_count + 1"),
		"last_used_at": now,
	}).Error
}

// FindAtomsBySynonym finds atoms where the search term matches value or synonyms
func (s *AtomService) FindAtomsBySynonym(searchTerm string) ([]models.Atom, error) {
	var atoms []models.Atom
	
	// Search in value, label, and synonyms
	searchPattern := "%" + searchTerm + "%"
	if err := s.db.Where("value LIKE ? OR label LIKE ? OR synonyms LIKE ?", 
		searchPattern, searchPattern, searchPattern).Find(&atoms).Error; err != nil {
		return nil, err
	}
	
	return atoms, nil
}

// GetPopularAtoms returns the most frequently used atoms
func (s *AtomService) GetPopularAtoms(limit int) ([]models.Atom, error) {
	var atoms []models.Atom
	if err := s.db.Order("usage_count DESC").Limit(limit).Find(&atoms).Error; err != nil {
		return nil, err
	}
	return atoms, nil
}

// BatchImportAtoms imports multiple atoms from JSON
func (s *AtomService) BatchImportAtoms(jsonData string) (int, error) {
	type ImportAtom struct {
		Value      string   `json:"value"`
		Label      string   `json:"label"`
		Type       string   `json:"type"`
		CategoryID uint     `json:"category_id"`
		Synonyms   []string `json:"synonyms"`
	}
	
	var atoms []ImportAtom
	if err := json.Unmarshal([]byte(jsonData), &atoms); err != nil {
		return 0, err
	}
	
	count := 0
	for _, a := range atoms {
		_, err := s.CreateAtom(a.Value, a.Label, a.Type, a.CategoryID, a.Synonyms)
		if err == nil {
			count++
		}
	}
	
	return count, nil
}
