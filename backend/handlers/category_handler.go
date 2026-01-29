package handlers

import (
	"promptmaster/backend/services"
)

// CategoryHandler exposes category operations to the frontend
type CategoryHandler struct {
	service *services.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CategoryResponse represents the response structure
type CategoryResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CreateCategoryRequest represents a create category request
type CreateCategoryRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	ParentID uint   `json:"parent_id"`
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(req CreateCategoryRequest) CategoryResponse {
	category, err := h.service.CreateCategory(req.Name, req.Type, req.ParentID)
	if err != nil {
		return CategoryResponse{Success: false, Error: err.Error()}
	}
	return CategoryResponse{Success: true, Data: category}
}

// GetCategoryByID retrieves a category by ID
func (h *CategoryHandler) GetCategoryByID(id uint) CategoryResponse {
	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		return CategoryResponse{Success: false, Error: err.Error()}
	}
	return CategoryResponse{Success: true, Data: category}
}

// GetCategoriesByParent retrieves categories by parent ID
func (h *CategoryHandler) GetCategoriesByParent(parentID uint, catType string) CategoryResponse {
	categories, err := h.service.GetCategoriesByParent(parentID, catType)
	if err != nil {
		return CategoryResponse{Success: false, Error: err.Error()}
	}
	return CategoryResponse{Success: true, Data: categories}
}

// GetCategoryTree retrieves the full category tree
func (h *CategoryHandler) GetCategoryTree(catType string) CategoryResponse {
	tree, err := h.service.GetCategoryTree(catType)
	if err != nil {
		return CategoryResponse{Success: false, Error: err.Error()}
	}
	return CategoryResponse{Success: true, Data: tree}
}

// UpdateCategoryRequest represents an update category request
type UpdateCategoryRequest struct {
	ID      uint                   `json:"id"`
	Updates map[string]interface{} `json:"updates"`
}

// UpdateCategory updates a category
func (h *CategoryHandler) UpdateCategory(req UpdateCategoryRequest) CategoryResponse {
	category, err := h.service.UpdateCategory(req.ID, req.Updates)
	if err != nil {
		return CategoryResponse{Success: false, Error: err.Error()}
	}
	return CategoryResponse{Success: true, Data: category}
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(id uint) CategoryResponse {
	if err := h.service.DeleteCategory(id); err != nil {
		return CategoryResponse{Success: false, Error: err.Error()}
	}
	return CategoryResponse{Success: true}
}

// MoveCategoryRequest represents a move category request
type MoveCategoryRequest struct {
	ID         uint `json:"id"`
	NewParentID uint `json:"new_parent_id"`
}

// MoveCategory moves a category to a new parent
func (h *CategoryHandler) MoveCategory(req MoveCategoryRequest) CategoryResponse {
	if err := h.service.MoveCategory(req.ID, req.NewParentID); err != nil {
		return CategoryResponse{Success: false, Error: err.Error()}
	}
	return CategoryResponse{Success: true}
}

// ReorderCategoriesRequest represents a reorder request
type ReorderCategoriesRequest struct {
	IDs []uint `json:"ids"`
}

// ReorderCategories reorders categories
func (h *CategoryHandler) ReorderCategories(req ReorderCategoriesRequest) CategoryResponse {
	if err := h.service.ReorderCategories(req.IDs); err != nil {
		return CategoryResponse{Success: false, Error: err.Error()}
	}
	return CategoryResponse{Success: true}
}
