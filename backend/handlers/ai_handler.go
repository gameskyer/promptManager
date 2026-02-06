package handlers

import (
	"promptmaster/backend/services"
)

// AIHandler exposes AI operations to the frontend
type AIHandler struct {
	service *services.AIService
}

// NewAIHandler creates a new AIHandler
func NewAIHandler(service *services.AIService) *AIHandler {
	return &AIHandler{service: service}
}

// AIResponse represents the response structure
type AIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ExplodePromptRequest represents a prompt explosion request
type ExplodePromptRequest struct {
	Prompt      string             `json:"prompt"`
	Categories  []string           `json:"categories,omitempty"`
	CategoryMap map[string]uint    `json:"category_map,omitempty"` // 分类名称到ID的映射
	Config      *services.AIConfig `json:"config,omitempty"`
}

// ExplodePrompt uses AI to break down a prompt
func (h *AIHandler) ExplodePrompt(req ExplodePromptRequest) AIResponse {
	result, err := h.service.ExplodePrompt(req.Prompt, req.Categories, req.CategoryMap, req.Config)
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true, Data: result}
}

// OptimizePromptRequest represents an optimize request
type OptimizePromptRequest struct {
	Prompt string             `json:"prompt"`
	Config *services.AIConfig `json:"config,omitempty"`
}

// OptimizePrompt uses AI to optimize a prompt
func (h *AIHandler) OptimizePrompt(req OptimizePromptRequest) AIResponse {
	result, err := h.service.OptimizePrompt(req.Prompt, req.Config)
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true, Data: result}
}

// TranslatePromptRequest represents a translate request
type TranslatePromptRequest struct {
	Prompt string             `json:"prompt"`
	Config *services.AIConfig `json:"config,omitempty"`
}

// TranslatePrompt uses AI to translate a prompt
func (h *AIHandler) TranslatePrompt(req TranslatePromptRequest) AIResponse {
	result, err := h.service.TranslatePrompt(req.Prompt, req.Config)
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true, Data: result}
}

// AnalyzePromptRequest represents an analyze request
type AnalyzePromptRequest struct {
	Prompt string             `json:"prompt"`
	Config *services.AIConfig `json:"config,omitempty"`
}

// AnalyzePrompt uses AI to analyze a prompt
func (h *AIHandler) AnalyzePrompt(req AnalyzePromptRequest) AIResponse {
	result, err := h.service.AnalyzePrompt(req.Prompt, req.Config)
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true, Data: result}
}

// ImportExtractedRequest represents an import request
type ImportExtractedRequest struct {
	Result     *services.ExplodeResult `json:"result"`
	CategoryID uint                    `json:"category_id"`
}

// ImportExtractedAtoms imports extracted atoms
func (h *AIHandler) ImportExtractedAtoms(req ImportExtractedRequest) AIResponse {
	importResult, err := h.service.ImportExtractedAtoms(req.Result, req.CategoryID)
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true, Data: importResult}
}

// ReverseImagePromptRequest represents an image-to-prompt request
type ReverseImagePromptRequest struct {
	ImagePath string `json:"image_path"`
}

// ReverseImagePrompt simulates image-to-prompt
func (h *AIHandler) ReverseImagePrompt(req ReverseImagePromptRequest) AIResponse {
	result, err := h.service.ReverseImagePrompt(req.ImagePath)
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true, Data: result}
}

// SaveAIConfigRequest represents a save config request
type SaveAIConfigRequest struct {
	Config *services.AIConfig `json:"config"`
}

// SaveAIConfig saves AI configuration
func (h *AIHandler) SaveAIConfig(req SaveAIConfigRequest) AIResponse {
	if err := h.service.SaveAIConfig(req.Config); err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true}
}

// GetAIConfig retrieves AI configuration
func (h *AIHandler) GetAIConfig() AIResponse {
	config, err := h.service.GetAIConfig()
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true, Data: config}
}

// GenericAIRequest represents a generic AI request for the unified endpoint
type GenericAIRequest struct {
	Mode        string             `json:"mode"` // explode, optimize, translate, analyze
	Prompt      string             `json:"prompt"`
	Categories  []string           `json:"categories,omitempty"`
	CategoryMap map[string]uint    `json:"category_map,omitempty"` // 分类名称到ID的映射
	Config      *services.AIConfig `json:"config,omitempty"`
}

// ProcessAI handles all AI operations through a unified endpoint
func (h *AIHandler) ProcessAI(req GenericAIRequest) AIResponse {
	switch req.Mode {
	case "explode":
		result, err := h.service.ExplodePrompt(req.Prompt, req.Categories, req.CategoryMap, req.Config)
		if err != nil {
			return AIResponse{Success: false, Error: err.Error()}
		}
		return AIResponse{Success: true, Data: result}
		
	case "optimize":
		result, err := h.service.OptimizePrompt(req.Prompt, req.Config)
		if err != nil {
			return AIResponse{Success: false, Error: err.Error()}
		}
		return AIResponse{Success: true, Data: result}
		
	case "translate":
		result, err := h.service.TranslatePrompt(req.Prompt, req.Config)
		if err != nil {
			return AIResponse{Success: false, Error: err.Error()}
		}
		return AIResponse{Success: true, Data: result}
		
	case "analyze":
		result, err := h.service.AnalyzePrompt(req.Prompt, req.Config)
		if err != nil {
			return AIResponse{Success: false, Error: err.Error()}
		}
		return AIResponse{Success: true, Data: result}
		
	default:
		return AIResponse{Success: false, Error: "unknown mode: " + req.Mode}
	}
}
