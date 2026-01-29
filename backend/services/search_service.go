package services

import (
	"fmt"
	"promptmaster/backend/models"
	"strings"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"gorm.io/gorm"
)

// SearchService handles search operations
type SearchService struct {
	db    *gorm.DB
	index bleve.Index
	mu    sync.RWMutex
}

// SearchResult represents a search result
type SearchResult struct {
	Type        string      `json:"type"` // atom, preset, category
	ID          uint        `json:"id"`
	Title       string      `json:"title"`
	Subtitle    string      `json:"subtitle"`
	MatchScore  float64     `json:"match_score"`
	Data        interface{} `json:"data,omitempty"`
}

// NewSearchService creates a new SearchService
func NewSearchService(db *gorm.DB) *SearchService {
	s := &SearchService{db: db}
	s.initIndex()
	return s
}

// initIndex initializes the Bleve search index
func (s *SearchService) initIndex() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Create an in-memory index
	mapping := bleve.NewIndexMapping()
	
	// Add custom analyzer for Chinese and Pinyin support
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = "standard"
	
	atomMapping := bleve.NewDocumentMapping()
	atomMapping.AddFieldMappingsAt("value", textFieldMapping)
	atomMapping.AddFieldMappingsAt("label", textFieldMapping)
	atomMapping.AddFieldMappingsAt("synonyms", textFieldMapping)
	
	mapping.AddDocumentMapping("atom", atomMapping)
	
	var err error
	s.index, err = bleve.NewMemOnly(mapping)
	if err != nil {
		fmt.Printf("Failed to create search index: %v\n", err)
		return
	}
	
	// Index existing atoms
	s.rebuildIndex()
}

// rebuildIndex rebuilds the search index from database
func (s *SearchService) rebuildIndex() {
	var atoms []models.Atom
	if err := s.db.Find(&atoms).Error; err != nil {
		return
	}
	
	for _, atom := range atoms {
		s.indexAtom(&atom)
	}
}

// indexAtom indexes a single atom
func (s *SearchService) indexAtom(atom *models.Atom) {
	doc := map[string]interface{}{
		"id":        atom.ID,
		"value":     atom.Value,
		"label":     atom.Label,
		"synonyms":  strings.Join(atom.Synonyms, " "),
		"type":      atom.Type,
		"category":  atom.CategoryID,
	}
	s.index.Index(fmt.Sprintf("atom_%d", atom.ID), doc)
}

// Search performs a global search
func (s *SearchService) Search(query string, limit int) ([]SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if limit <= 0 {
		limit = 20
	}
	
	// Search in Bleve index
	bleveQuery := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(bleveQuery)
	searchRequest.Size = limit
	searchRequest.Fields = []string{"*"}
	
	searchResults, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	
	var results []SearchResult
	for _, hit := range searchResults.Hits {
		idStr := strings.TrimPrefix(hit.ID, "atom_")
		var atomID uint
		fmt.Sscanf(idStr, "%d", &atomID)
		
		var atom models.Atom
		if err := s.db.Preload("Category").First(&atom, atomID).Error; err != nil {
			continue
		}
		
		results = append(results, SearchResult{
			Type:       "atom",
			ID:         atom.ID,
			Title:      atom.Label,
			Subtitle:   atom.Value,
			MatchScore: hit.Score,
			Data:       atom,
		})
	}
	
	return results, nil
}

// SearchAtoms searches only for atoms
func (s *SearchService) SearchAtoms(searchTerm string, atomType string, categoryID uint, limit int) ([]models.Atom, error) {
	var atoms []models.Atom
	
	query := s.db.Model(&models.Atom{})
	
	if searchTerm != "" {
		searchPattern := "%" + searchTerm + "%"
		query = query.Where("value LIKE ? OR label LIKE ? OR synonyms LIKE ?", 
			searchPattern, searchPattern, searchPattern)
	}
	
	if atomType != "" {
		query = query.Where("type = ?", atomType)
	}
	
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Order("usage_count DESC").Find(&atoms).Error; err != nil {
		return nil, err
	}
	
	return atoms, nil
}

// QuickSearch performs a quick search with pinyin support
func (s *SearchService) QuickSearch(term string) ([]SearchResult, error) {
	results, err := s.Search(term, 10)
	if err != nil {
		return nil, err
	}
	
	// Also search by pinyin initials
	pinyinResults := s.searchByPinyin(term)
	
	// Merge results
	seen := make(map[uint]bool)
	var merged []SearchResult
	
	for _, r := range results {
		if !seen[r.ID] {
			seen[r.ID] = true
			merged = append(merged, r)
		}
	}
	
	for _, r := range pinyinResults {
		if !seen[r.ID] {
			seen[r.ID] = true
			merged = append(merged, r)
		}
	}
	
	return merged, nil
}

// searchByPinyin searches atoms by pinyin initials
func (s *SearchService) searchByPinyin(pinyin string) []SearchResult {
	// This is a simplified implementation
	// In production, you'd use a proper pinyin library
	var results []SearchResult
	var atoms []models.Atom
	
	// Get all atoms and do simple pinyin matching
	s.db.Find(&atoms)
	
	for _, atom := range atoms {
		// Simple pinyin matching - check if pinyin initials match
		// For example: "服装" (fuzhuang) should match "fz"
		if matchPinyin(atom.Label, pinyin) {
			results = append(results, SearchResult{
				Type:     "atom",
				ID:       atom.ID,
				Title:    atom.Label,
				Subtitle: atom.Value,
				Data:     atom,
			})
		}
	}
	
	return results
}

// matchPinyin checks if pinyin initials match the search term
func matchPinyin(chinese, pinyin string) bool {
	// Simplified implementation - in production use go-pinyin library
	// For now, just do a simple contains check
	return strings.Contains(strings.ToLower(chinese), strings.ToLower(pinyin))
}

// ReindexAtom updates the index for a specific atom
func (s *SearchService) ReindexAtom(atom *models.Atom) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.indexAtom(atom)
}

// RemoveFromIndex removes an atom from the index
func (s *SearchService) RemoveFromIndex(atomID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.index.Delete(fmt.Sprintf("atom_%d", atomID))
}

// SearchPresets searches for presets by title
func (s *SearchService) SearchPresets(searchTerm string, limit int) ([]models.Preset, error) {
	var presets []models.Preset
	
	searchPattern := "%" + searchTerm + "%"
	query := s.db.Where("title LIKE ? AND is_deleted = ?", searchPattern, false)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Order("updated_at DESC").Find(&presets).Error; err != nil {
		return nil, err
	}
	
	return presets, nil
}
