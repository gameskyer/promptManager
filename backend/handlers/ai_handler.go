package handlers

import (
	"promptmaster/backend/models"
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
	Prompt             string             `json:"prompt"`
	Categories         []string           `json:"categories,omitempty"`
	CategoryMap        map[string]uint    `json:"category_map,omitempty"` // 分类名称到ID的映射
	Config             *services.AIConfig `json:"config,omitempty"`
	SystemPrompt       string             `json:"system_prompt,omitempty"`        // 自定义系统提示词模板
	UserPromptTemplate string             `json:"user_prompt_template,omitempty"` // 用户提示词模板
}

// ExplodePrompt uses AI to break down a prompt
func (h *AIHandler) ExplodePrompt(req ExplodePromptRequest) AIResponse {
	result, err := h.service.ExplodePrompt(req.Prompt, req.Categories, req.CategoryMap, req.Config, req.SystemPrompt, req.UserPromptTemplate)
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true, Data: result}
}

// OptimizePromptRequest represents an optimize request
type OptimizePromptRequest struct {
	Prompt             string             `json:"prompt"`
	Config             *services.AIConfig `json:"config,omitempty"`
	SystemPrompt       string             `json:"system_prompt,omitempty"`        // 自定义系统提示词模板
	UserPromptTemplate string             `json:"user_prompt_template,omitempty"` // 用户提示词模板
}

// OptimizePrompt uses AI to optimize a prompt
func (h *AIHandler) OptimizePrompt(req OptimizePromptRequest) AIResponse {
	result, err := h.service.OptimizePrompt(req.Prompt, req.Config, req.SystemPrompt, req.UserPromptTemplate)
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true, Data: result}
}

// TranslatePromptRequest represents a translate request
type TranslatePromptRequest struct {
	Prompt             string             `json:"prompt"`
	Config             *services.AIConfig `json:"config,omitempty"`
	SystemPrompt       string             `json:"system_prompt,omitempty"`        // 自定义系统提示词模板
	UserPromptTemplate string             `json:"user_prompt_template,omitempty"` // 用户提示词模板
}

// TranslatePrompt uses AI to translate a prompt
func (h *AIHandler) TranslatePrompt(req TranslatePromptRequest) AIResponse {
	result, err := h.service.TranslatePrompt(req.Prompt, req.Config, req.SystemPrompt, req.UserPromptTemplate)
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true, Data: result}
}

// AnalyzePromptRequest represents an analyze request
type AnalyzePromptRequest struct {
	Prompt             string             `json:"prompt"`
	Config             *services.AIConfig `json:"config,omitempty"`
	SystemPrompt       string             `json:"system_prompt,omitempty"`        // 自定义系统提示词模板
	UserPromptTemplate string             `json:"user_prompt_template,omitempty"` // 用户提示词模板
}

// AnalyzePrompt uses AI to analyze a prompt
func (h *AIHandler) AnalyzePrompt(req AnalyzePromptRequest) AIResponse {
	result, err := h.service.AnalyzePrompt(req.Prompt, req.Config, req.SystemPrompt, req.UserPromptTemplate)
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
	Mode               string             `json:"mode"` // explode, optimize, translate, analyze
	Prompt             string             `json:"prompt"`
	Categories         []string           `json:"categories,omitempty"`
	CategoryMap        map[string]uint    `json:"category_map,omitempty"` // 分类名称到ID的映射
	Config             *services.AIConfig `json:"config,omitempty"`
	SystemPrompt       string             `json:"system_prompt,omitempty"`        // 自定义系统提示词模板
	UserPromptTemplate string             `json:"user_prompt_template,omitempty"` // 用户提示词模板
}

// ProcessAI handles all AI operations through a unified endpoint
func (h *AIHandler) ProcessAI(req GenericAIRequest) AIResponse {
	switch req.Mode {
	case "explode":
		result, err := h.service.ExplodePrompt(req.Prompt, req.Categories, req.CategoryMap, req.Config, req.SystemPrompt, req.UserPromptTemplate)
		if err != nil {
			return AIResponse{Success: false, Error: err.Error()}
		}
		return AIResponse{Success: true, Data: result}

	case "optimize":
		result, err := h.service.OptimizePrompt(req.Prompt, req.Config, req.SystemPrompt, req.UserPromptTemplate)
		if err != nil {
			return AIResponse{Success: false, Error: err.Error()}
		}
		return AIResponse{Success: true, Data: result}

	case "translate":
		result, err := h.service.TranslatePrompt(req.Prompt, req.Config, req.SystemPrompt, req.UserPromptTemplate)
		if err != nil {
			return AIResponse{Success: false, Error: err.Error()}
		}
		return AIResponse{Success: true, Data: result}

	case "analyze":
		result, err := h.service.AnalyzePrompt(req.Prompt, req.Config, req.SystemPrompt, req.UserPromptTemplate)
		if err != nil {
			return AIResponse{Success: false, Error: err.Error()}
		}
		return AIResponse{Success: true, Data: result}

	default:
		return AIResponse{Success: false, Error: "unknown mode: " + req.Mode}
	}
}

// ProviderConfigDTO represents AI provider config for API response
type ProviderConfigDTO struct {
	Provider string                 `json:"provider"`
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	BaseURL  string                 `json:"base_url"`
	APIKey   string                 `json:"api_key"`
	Model    string                 `json:"model"`
	Models   []string               `json:"models"`
	Headers  map[string]interface{} `json:"headers"`
	Enabled  bool                   `json:"enabled"`
	IsCustom bool                   `json:"is_custom"`
}

// PromptTemplateDTO represents AI prompt template for API response
type PromptTemplateDTO struct {
	TemplateID       string  `json:"template_id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	SystemPrompt     string  `json:"system_prompt"`
	UserPromptTemplate string `json:"user_prompt_template"`
	Temperature      float64 `json:"temperature"`
	ResponseFormat   string  `json:"response_format"`
	IsCustom         bool    `json:"is_custom"`
}

// SettingsDTO represents AI settings for API response
type SettingsDTO struct {
	CurrentProvider string `json:"current_provider"`
	CurrentPrompt   string `json:"current_prompt"`
}

// ============== Provider Configuration APIs ==============

// GetAllProviders retrieves all AI provider configurations
func (h *AIHandler) GetAllProviders() AIResponse {
	providers, err := h.service.GetAllProviders()
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	
	// Convert to DTO
	var dtos []ProviderConfigDTO
	for _, p := range providers {
		dtos = append(dtos, ProviderConfigDTO{
			Provider: p.Provider,
			Name:     p.Name,
			Type:     p.Type,
			BaseURL:  p.BaseURL,
			APIKey:   p.APIKey,
			Model:    p.Model,
			Models:   p.Models,
			Headers:  p.Headers,
			Enabled:  p.Enabled,
			IsCustom: p.IsCustom,
		})
	}
	
	return AIResponse{Success: true, Data: dtos}
}

// SaveProviderRequest represents a save provider request
type SaveProviderRequest struct {
	Provider ProviderConfigDTO `json:"provider"`
}

// SaveProvider saves an AI provider configuration
func (h *AIHandler) SaveProvider(req SaveProviderRequest) AIResponse {
	// Convert DTO to model
	provider := models.AIProviderConfig{
		Provider: req.Provider.Provider,
		Name:     req.Provider.Name,
		Type:     req.Provider.Type,
		BaseURL:  req.Provider.BaseURL,
		APIKey:   req.Provider.APIKey,
		Model:    req.Provider.Model,
		Models:   req.Provider.Models,
		Headers:  req.Provider.Headers,
		Enabled:  req.Provider.Enabled,
		IsCustom: req.Provider.IsCustom,
	}
	
	if err := h.service.SaveProvider(&provider); err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true}
}

// DeleteProvider deletes an AI provider configuration
func (h *AIHandler) DeleteProvider(providerID string) AIResponse {
	if err := h.service.DeleteProvider(providerID); err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true}
}

// SetCurrentProviderRequest represents set current provider request
type SetCurrentProviderRequest struct {
	ProviderID string `json:"provider_id"`
}

// SetCurrentProvider sets the current active provider
func (h *AIHandler) SetCurrentProvider(req SetCurrentProviderRequest) AIResponse {
	if err := h.service.SetCurrentProvider(req.ProviderID); err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true}
}

// ============== Prompt Template APIs ==============

// GetAllPromptTemplates retrieves all prompt templates
func (h *AIHandler) GetAllPromptTemplates() AIResponse {
	templates, err := h.service.GetAllPromptTemplates()
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	
	// Convert to DTO
	var dtos []PromptTemplateDTO
	for _, t := range templates {
		dtos = append(dtos, PromptTemplateDTO{
			TemplateID:       t.TemplateID,
			Name:             t.Name,
			Description:      t.Description,
			SystemPrompt:     t.SystemPrompt,
			UserPromptTemplate: t.UserPromptTemplate,
			Temperature:      t.Temperature,
			ResponseFormat:   t.ResponseFormat,
			IsCustom:         t.IsCustom,
		})
	}
	
	return AIResponse{Success: true, Data: dtos}
}

// GetPromptTemplate retrieves a prompt template by ID
func (h *AIHandler) GetPromptTemplate(templateID string) AIResponse {
	template, err := h.service.GetPromptTemplate(templateID)
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	
	dto := PromptTemplateDTO{
		TemplateID:       template.TemplateID,
		Name:             template.Name,
		Description:      template.Description,
		SystemPrompt:     template.SystemPrompt,
		UserPromptTemplate: template.UserPromptTemplate,
		Temperature:      template.Temperature,
		ResponseFormat:   template.ResponseFormat,
		IsCustom:         template.IsCustom,
	}
	
	return AIResponse{Success: true, Data: dto}
}

// SavePromptTemplateRequest represents save prompt template request
type SavePromptTemplateRequest struct {
	Template PromptTemplateDTO `json:"template"`
}

// SavePromptTemplate saves a prompt template
func (h *AIHandler) SavePromptTemplate(req SavePromptTemplateRequest) AIResponse {
	// Convert DTO to model
	template := models.AIPromptTemplate{
		TemplateID:       req.Template.TemplateID,
		Name:             req.Template.Name,
		Description:      req.Template.Description,
		SystemPrompt:     req.Template.SystemPrompt,
		UserPromptTemplate: req.Template.UserPromptTemplate,
		Temperature:      req.Template.Temperature,
		ResponseFormat:   req.Template.ResponseFormat,
		IsCustom:         req.Template.IsCustom,
	}
	
	if err := h.service.SavePromptTemplate(&template); err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true}
}

// DeletePromptTemplate deletes a prompt template
func (h *AIHandler) DeletePromptTemplate(templateID string) AIResponse {
	if err := h.service.DeletePromptTemplate(templateID); err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true}
}

// SetCurrentPromptRequest represents set current prompt request
type SetCurrentPromptRequest struct {
	TemplateID string `json:"template_id"`
}

// SetCurrentPrompt sets the current active prompt template
func (h *AIHandler) SetCurrentPrompt(req SetCurrentPromptRequest) AIResponse {
	if err := h.service.SetCurrentPrompt(req.TemplateID); err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	return AIResponse{Success: true}
}

// GetCurrentSettings retrieves current AI settings
func (h *AIHandler) GetCurrentSettings() AIResponse {
	settings, err := h.service.GetCurrentSettings()
	if err != nil {
		return AIResponse{Success: false, Error: err.Error()}
	}
	
	dto := SettingsDTO{
		CurrentProvider: settings.CurrentProvider,
		CurrentPrompt:   settings.CurrentPrompt,
	}
	
	return AIResponse{Success: true, Data: dto}
}
