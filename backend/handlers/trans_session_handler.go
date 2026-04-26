package handlers

import (
	"promptmaster/backend/services"
)

// TransSessionHandler exposes translation session operations to the frontend
type TransSessionHandler struct {
	service *services.TransSessionService
}

// NewTransSessionHandler creates a new TransSessionHandler
func NewTransSessionHandler(service *services.TransSessionService) *TransSessionHandler {
	return &TransSessionHandler{service: service}
}

// TransSessionResponse represents the response structure
type TransSessionResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SplitPromptRequest represents a split request with direction
type SplitPromptRequest struct {
	Text      string `json:"text"`
	Direction string `json:"direction"` // "cn2en" or "en2cn"
}

// SplitPrompt splits raw prompt text into individual words
func (h *TransSessionHandler) SplitPrompt(req SplitPromptRequest) TransSessionResponse {
	words := h.service.SplitPrompt(req.Text, services.TranslationDirection(req.Direction))
	return TransSessionResponse{Success: true, Data: words}
}

// TranslateWordsRequest represents a batch translate request
type TranslateWordsRequest struct {
	Words     []services.PromptWord `json:"words"`
	Direction string                `json:"direction"` // "cn2en" or "en2cn"
}

// TranslateWords translates all untranslated words in the list
func (h *TransSessionHandler) TranslateWords(req TranslateWordsRequest) TransSessionResponse {
	words, err := h.service.TranslateWords(req.Words, services.TranslationDirection(req.Direction))
	if err != nil {
		return TransSessionResponse{Success: false, Error: err.Error()}
	}
	return TransSessionResponse{Success: true, Data: words}
}

// TranslateSingleRequest represents a single word translate request
type TranslateSingleRequest struct {
	Text      string `json:"text"`
	Direction string `json:"direction"` // "cn2en" or "en2cn"
}

// TranslateSingle translates a single word
func (h *TransSessionHandler) TranslateSingle(req TranslateSingleRequest) TransSessionResponse {
	result, err := h.service.TranslateSingle(req.Text, services.TranslationDirection(req.Direction))
	if err != nil {
		return TransSessionResponse{Success: false, Error: err.Error()}
	}
	return TransSessionResponse{Success: true, Data: result}
}
