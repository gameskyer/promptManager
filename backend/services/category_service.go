package services

import (
	"promptmaster/backend/models"

	"gorm.io/gorm"
)

// CategoryService handles category-related operations
type CategoryService struct {
	db *gorm.DB
}

// NewCategoryService creates a new CategoryService
func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(name, catType string, parentID uint) (*models.Category, error) {
	// Get max sort order
	var maxSort int
	s.db.Model(&models.Category{}).Where("parent_id = ?", parentID).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxSort)
	
	category := &models.Category{
		Name:      name,
		Type:      catType,
		ParentID:  parentID,
		SortOrder: maxSort + 1,
	}
	
	if err := s.db.Create(category).Error; err != nil {
		return nil, err
	}
	
	return category, nil
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetCategoriesByParent retrieves categories by parent ID
func (s *CategoryService) GetCategoriesByParent(parentID uint, catType string) ([]models.Category, error) {
	var categories []models.Category
	query := s.db.Where("parent_id = ?", parentID).Order("sort_order ASC")
	if catType != "" {
		query = query.Where("type = ?", catType)
	}
	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryTree returns the full category tree
func (s *CategoryService) GetCategoryTree(catType string) ([]CategoryNode, error) {
	var categories []models.Category
	query := s.db.Order("sort_order ASC")
	if catType != "" {
		query = query.Where("type = ?", catType)
	}
	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}
	
	return buildTree(categories, 0), nil
}

// CategoryNode represents a node in the category tree
type CategoryNode struct {
	models.Category
	Children []CategoryNode `json:"children"`
}

func buildTree(categories []models.Category, parentID uint) []CategoryNode {
	var nodes []CategoryNode
	for _, cat := range categories {
		if cat.ParentID == parentID {
			node := CategoryNode{
				Category: cat,
				Children: buildTree(categories, cat.ID),
			}
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(id uint, updates map[string]interface{}) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	
	if err := s.db.Model(&category).Updates(updates).Error; err != nil {
		return nil, err
	}
	
	return &category, nil
}

// DeleteCategory deletes a category and its children
func (s *CategoryService) DeleteCategory(id uint) error {
	// Delete children first
	var children []models.Category
	s.db.Where("parent_id = ?", id).Find(&children)
	for _, child := range children {
		s.DeleteCategory(child.ID)
	}
	
	return s.db.Delete(&models.Category{}, id).Error
}

// MoveCategory moves a category to a new parent
func (s *CategoryService) MoveCategory(id, newParentID uint) error {
	return s.db.Model(&models.Category{}).Where("id = ?", id).Update("parent_id", newParentID).Error
}

// ReorderCategories updates the sort order of categories
func (s *CategoryService) ReorderCategories(ids []uint) error {
	for i, id := range ids {
		if err := s.db.Model(&models.Category{}).Where("id = ?", id).Update("sort_order", i+1).Error; err != nil {
			return err
		}
	}
	return nil
}
