package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"promptmaster/backend/logger"

	"gorm.io/gorm"
)

// AIService handles AI-powered operations
type AIService struct {
	db          *gorm.DB
	atomService *AtomService
	httpClient  *http.Client
	logger      *logger.Logger
}

// AIConfig holds AI service configuration
type AIConfig struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
	Endpoint string `json:"endpoint"`
	Model    string `json:"model"`
}

// ExplodeResult represents the result of AI text explosion
type ExplodeResult struct {
	Atoms     []ExtractedAtom `json:"atoms"`
	RawPrompt string          `json:"raw_prompt"`
}

// ExtractedAtom represents an atom extracted from text
type ExtractedAtom struct {
	Value        string   `json:"value"`
	Label        string   `json:"label"`
	Type         string   `json:"type"`
	CategoryID   uint     `json:"category"`   // 分类ID
	CategoryName string   `json:"category_name,omitempty"` // 分类名称（用于展示）
	Synonyms     []string `json:"synonyms"`
	IsNew        bool     `json:"is_new"`
	ExistingID   uint     `json:"existing_id,omitempty"`
}

// OptimizeResult represents the result of prompt optimization
type OptimizeResult struct {
	Optimized   string   `json:"optimized"`
	Changes     []string `json:"changes"`
	Suggestions []string `json:"suggestions"`
}

// TranslateResult represents the result of prompt translation
type TranslateResult struct {
	Translation string   `json:"translation"`
	Keywords    []string `json:"keywords"`
	Notes       string   `json:"notes"`
}

// AnalyzeResult represents the result of prompt analysis
type AnalyzeResult struct {
	Analysis    map[string]string `json:"analysis"`
	Issues      []string          `json:"issues"`
	Suggestions []string          `json:"suggestions"`
}

// NewAIService creates a new AIService
func NewAIService(db *gorm.DB) *AIService {
	// 初始化日志记录器（10MB分割）
	log, err := logger.NewLogger(logger.GetAILogPath(), 10)
	if err != nil {
		fmt.Printf("Warning: failed to create AI logger: %v\n", err)
	}

	return &AIService{
		db:          db,
		atomService: NewAtomService(db),
		httpClient:  &http.Client{Timeout: 180 * time.Second},
		logger:      log,
	}
}

// ExplodePrompt uses AI to break down a long prompt into atomic words
func (s *AIService) ExplodePrompt(prompt string, categories []string, categoryMap map[string]uint, config *AIConfig) (*ExplodeResult, error) {
	if config == nil || config.APIKey == "" {
		return s.ruleBasedExplosion(prompt, categoryMap)
	}
	return s.aiBasedExplosion(prompt, categories, categoryMap, config)
}

// OptimizePrompt uses AI to optimize a prompt
func (s *AIService) OptimizePrompt(prompt string, config *AIConfig) (*OptimizeResult, error) {
	if config == nil || config.APIKey == "" {
		return s.ruleBasedOptimization(prompt)
	}
	return s.aiBasedOptimization(prompt, config)
}

// TranslatePrompt uses AI to translate a prompt
func (s *AIService) TranslatePrompt(prompt string, config *AIConfig) (*TranslateResult, error) {
	if config == nil || config.APIKey == "" {
		return s.ruleBasedTranslation(prompt)
	}
	return s.aiBasedTranslation(prompt, config)
}

// AnalyzePrompt uses AI to analyze a prompt
func (s *AIService) AnalyzePrompt(prompt string, config *AIConfig) (*AnalyzeResult, error) {
	if config == nil || config.APIKey == "" {
		return s.ruleBasedAnalysis(prompt)
	}
	return s.aiBasedAnalysis(prompt, config)
}

// ruleBasedExplosion extracts atoms using rule-based parsing
func (s *AIService) ruleBasedExplosion(prompt string, categoryMap map[string]uint) (*ExplodeResult, error) {
	delimiters := regexp.MustCompile(`[,，;；|\/]`)
	parts := delimiters.Split(prompt, -1)

	var atoms []ExtractedAtom

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		existingAtom, err := s.atomService.GetAtomByValue(strings.ToLower(part))

		atom := ExtractedAtom{
			Value:      part,
			Label:      part,
			Type:       "Positive",
			IsNew:      err != nil || existingAtom == nil,
			CategoryID: 0, // 默认无分类
		}

		if err == nil && existingAtom != nil {
			atom.ExistingID = existingAtom.ID
			atom.Label = existingAtom.Label
			atom.Type = existingAtom.Type
			atom.Synonyms = existingAtom.Synonyms
			atom.CategoryID = existingAtom.CategoryID
			atom.CategoryName = existingAtom.Category.Name
		}

		atoms = append(atoms, atom)
	}

	return &ExplodeResult{
		Atoms:     atoms,
		RawPrompt: prompt,
	}, nil
}

// aiBasedExplosion uses AI API to extract atoms
func (s *AIService) aiBasedExplosion(prompt string, categories []string, categoryMap map[string]uint, config *AIConfig) (*ExplodeResult, error) {
	// 如果 categoryMap 为空，使用默认分类
	if len(categoryMap) == 0 {
		categoryMap = map[string]uint{
			"质量":    1,
			"人物":    2,
			"姿势":    3,
			"场景":    4,
			"服装":    5,
			"发型":    6,
			"道具":    7,
			"风格":    8,
			"光照":    9,
			"其他":    10,
		}
	}
	
	// 构建分类列表（名称:ID格式）
	var categoryEntries []string
	for name, id := range categoryMap {
		categoryEntries = append(categoryEntries, fmt.Sprintf("%s(%d)", name, id))
	}
	categoryList := strings.Join(categoryEntries, ", ")
	
	// 构建分类名称列表（用于AI理解）
	var categoryNames []string
	for name := range categoryMap {
		categoryNames = append(categoryNames, name)
	}
	categoryNameList := strings.Join(categoryNames, ", ")

	systemPrompt := fmt.Sprintf(`You are a prompt engineering expert for AI image generation. Break down the following prompt into atomic components.

Available categories with their IDs:
%s

Rules:
1. Each component should be a single concept or attribute
2. Use English for value field (lowercase, underscore for multi-word)
3. Provide Chinese translation in label field
4. Classify as Positive or Negative type
5. Assign category ID (number) based on the available categories above. Choose the most appropriate category ID from: %s
6. Include relevant synonyms
7. Return category as a NUMBER (ID), NOT as a string

Return JSON format strictly:
{"atoms": [{"value": "masterpiece", "label": "杰作", "type": "Positive", "category": 1, "synonyms": ["best quality"]}]}`, categoryList, categoryNameList)

	response, err := s.callAIAPI(config, systemPrompt, prompt)
	if err != nil {
		return nil, err
	}

	var result ExplodeResult
	if err := s.extractJSON(response, &result); err != nil {
		var atoms []ExtractedAtom
		if err := s.extractJSON(response, &atoms); err != nil {
			return nil, fmt.Errorf("failed to parse AI response: %w", err)
		}
		result.Atoms = atoms
	}

	result.RawPrompt = prompt

	for i := range result.Atoms {
		existing, err := s.atomService.GetAtomByValue(strings.ToLower(result.Atoms[i].Value))
		if err == nil && existing != nil {
			result.Atoms[i].IsNew = false
			result.Atoms[i].ExistingID = existing.ID
			// 如果数据库中有分类，优先使用数据库的分类信息
			if existing.CategoryID > 0 {
				result.Atoms[i].CategoryID = existing.CategoryID
				result.Atoms[i].CategoryName = existing.Category.Name
			}
		} else {
			result.Atoms[i].IsNew = true
		}
		
		// 根据CategoryID查找分类名称（对于新词或数据库中没有分类的词）
		if result.Atoms[i].CategoryName == "" && result.Atoms[i].CategoryID > 0 {
			for name, id := range categoryMap {
				if id == result.Atoms[i].CategoryID {
					result.Atoms[i].CategoryName = name
					break
				}
			}
		}
		
		// 如果仍然没有分类名称，标记为未分类
		if result.Atoms[i].CategoryName == "" {
			result.Atoms[i].CategoryName = "未分类"
		}
	}

	return &result, nil
}

// aiBasedOptimization uses AI to optimize prompt
func (s *AIService) aiBasedOptimization(prompt string, config *AIConfig) (*OptimizeResult, error) {
	systemPrompt := "You are an AI image generation prompt optimization expert. Optimize the given prompt for better quality and results. Return JSON format strictly: {\"optimized\": \"optimized prompt\", \"changes\": [\"change 1\"], \"suggestions\": [\"suggestion 1\"]}"

	response, err := s.callAIAPI(config, systemPrompt, prompt)
	if err != nil {
		return nil, err
	}

	var result OptimizeResult
	if err := s.extractJSON(response, &result); err != nil {
		return &OptimizeResult{
			Optimized:   s.cleanAIResponse(response),
			Changes:     []string{"AI optimized"},
			Suggestions: []string{},
		}, nil
	}

	return &result, nil
}

// aiBasedTranslation uses AI to translate prompt
func (s *AIService) aiBasedTranslation(prompt string, config *AIConfig) (*TranslateResult, error) {
	systemPrompt := "You are a professional translator for AI image generation prompts. Translate the given Chinese prompt into high-quality English. Return JSON format strictly: {\"translation\": \"English text\", \"keywords\": [\"keyword1\"], \"notes\": \"notes\"}"

	response, err := s.callAIAPI(config, systemPrompt, prompt)
	if err != nil {
		return nil, err
	}

	var result TranslateResult
	if err := s.extractJSON(response, &result); err != nil {
		return &TranslateResult{
			Translation: s.cleanAIResponse(response),
			Keywords:    []string{},
			Notes:       "Translated by AI",
		}, nil
	}

	return &result, nil
}

// aiBasedAnalysis uses AI to analyze prompt
func (s *AIService) aiBasedAnalysis(prompt string, config *AIConfig) (*AnalyzeResult, error) {
	systemPrompt := "You are an AI image generation prompt analysis expert. Analyze the given prompt. Return JSON format strictly: {\"analysis\": {\"subject\": \"desc\", \"style\": \"desc\", \"quality\": \"desc\", \"lighting\": \"desc\", \"other\": \"desc\"}, \"issues\": [], \"suggestions\": []}"

	response, err := s.callAIAPI(config, systemPrompt, prompt)
	if err != nil {
		return nil, err
	}

	var result AnalyzeResult
	if err := s.extractJSON(response, &result); err != nil {
		return &AnalyzeResult{
			Analysis: map[string]string{
				"subject":  "AI analysis needed",
				"style":    "AI analysis needed",
				"quality":  "AI analysis needed",
				"lighting": "AI analysis needed",
				"other":    "AI analysis needed",
			},
			Issues:      []string{"Parse error"},
			Suggestions: []string{"Please retry"},
		}, nil
	}

	return &result, nil
}

// callAIAPI makes the actual API call to AI service
func (s *AIService) callAIAPI(config *AIConfig, systemPrompt, userPrompt string) (string, error) {
	startTime := time.Now()

	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1"
	}

	model := config.Model
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	requestBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.3,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 记录请求日志
	if s.logger != nil {
		s.logger.LogAIRequest(config.Provider, endpoint+"/chat/completions", model, requestBody)
	}

	req, err := http.NewRequest("POST", endpoint+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		if s.logger != nil {
			s.logger.Error(fmt.Sprintf("Failed to create request: %v", err))
		}
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(fmt.Sprintf("Request failed: %v", err))
		}
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(fmt.Sprintf("Failed to read response: %v", err))
		}
		return "", err
	}
	bodyString := string(bodyBytes)

	// 记录响应日志
	duration := time.Since(startTime)
	if s.logger != nil {
		s.logger.LogAIResponse(resp.StatusCode, bodyString, duration)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.Unmarshal(bodyBytes, &errResp)
		return "", fmt.Errorf("AI API error (%d): %v", resp.StatusCode, errResp)
	}

	var aiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &aiResponse); err != nil {
		return "", err
	}

	if len(aiResponse.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return aiResponse.Choices[0].Message.Content, nil
}

// extractJSON extracts JSON from AI response
func (s *AIService) extractJSON(content string, v interface{}) error {
	// Try code blocks
	re := regexp.MustCompile("`{3}(?:json)?\\s*([\\s\\S]*?)`{3}")
	matches := re.FindStringSubmatch(content)
	s.logger.Info("AI Response:" + content)

	var jsonStr string
	if len(matches) > 1 {
		jsonStr = matches[1]
	} else {
		// Find between braces
		start := strings.Index(content, "{")
		end := strings.LastIndex(content, "}")
		if start != -1 && end != -1 && end > start {
			jsonStr = content[start : end+1]
		} else {
			// Array format
			start = strings.Index(content, "[")
			end = strings.LastIndex(content, "]")
			if start != -1 && end != -1 && end > start {
				jsonStr = content[start : end+1]
			} else {
				jsonStr = content
			}
		}
	}
	return json.Unmarshal([]byte(jsonStr), v)
}

// cleanAIResponse removes markdown
func (s *AIService) cleanAIResponse(response string) string {
	re := regexp.MustCompile("`{3}[\\w]*\\n?")
	cleaned := re.ReplaceAllString(response, "")
	cleaned = strings.ReplaceAll(cleaned, "```", "")
	return strings.TrimSpace(cleaned)
}

// Rule-based fallbacks
func (s *AIService) ruleBasedOptimization(prompt string) (*OptimizeResult, error) {
	terms := strings.Split(prompt, ",")
	seen := make(map[string]bool)
	var unique []string

	for _, term := range terms {
		term = strings.TrimSpace(strings.ToLower(term))
		if term != "" && !seen[term] {
			seen[term] = true
			unique = append(unique, term)
		}
	}

	return &OptimizeResult{
		Optimized:   strings.Join(unique, ", "),
		Changes:     []string{"Removed duplicates"},
		Suggestions: []string{"Add quality tags"},
	}, nil
}

func (s *AIService) ruleBasedTranslation(prompt string) (*TranslateResult, error) {
	return &TranslateResult{
		Translation: prompt,
		Keywords:    strings.Split(prompt, ","),
		Notes:       "Local mode - configure AI API for better results",
	}, nil
}

func (s *AIService) ruleBasedAnalysis(prompt string) (*AnalyzeResult, error) {
	terms := strings.Split(prompt, ",")
	issues := []string{}

	if len(terms) < 5 {
		issues = append(issues, "Too few tags")
	}

	return &AnalyzeResult{
		Analysis: map[string]string{
			"subject":  "Rule-based analysis",
			"style":    "Need AI analysis",
			"quality":  fmt.Sprintf("%d tags", len(terms)),
			"lighting": "Need AI analysis",
			"other":    "Local mode",
		},
		Issues:      issues,
		Suggestions: []string{"Configure AI API"},
	}, nil
}

// ImportExtractedAtoms imports extracted atoms into the database
func (s *AIService) ImportExtractedAtoms(result *ExplodeResult, defaultCategoryID uint) (*ImportResult, error) {
	imported := 0
	updated := 0

	for _, atom := range result.Atoms {
		// 使用原子词自己的CategoryID，如果没有则使用默认分类ID
		atomCategoryID := atom.CategoryID
		if atomCategoryID == 0 {
			atomCategoryID = defaultCategoryID
		}
		
		if atom.IsNew {
			_, err := s.atomService.CreateAtom(
				atom.Value,
				atom.Label,
				atom.Type,
				atomCategoryID,
				atom.Synonyms,
			)
			if err == nil {
				imported++
			}
		} else if atom.ExistingID > 0 {
			if len(atom.Synonyms) > 0 {
				existing, _ := s.atomService.GetAtomByID(atom.ExistingID)
				if existing != nil {
					synonymMap := make(map[string]bool)
					for _, s := range existing.Synonyms {
						synonymMap[s] = true
					}
					for _, s := range atom.Synonyms {
						synonymMap[s] = true
					}

					var mergedSynonyms []string
					for s := range synonymMap {
						mergedSynonyms = append(mergedSynonyms, s)
					}

					s.atomService.UpdateAtom(atom.ExistingID, map[string]interface{}{
						"synonyms": mergedSynonyms,
					})
					updated++
				}
			}
		}
	}

	return &ImportResult{
		Imported: imported,
		Updated:  updated,
	}, nil
}

// ImportResult represents the result of importing atoms
type ImportResult struct {
	Imported int `json:"imported"`
	Updated  int `json:"updated"`
}

// ReverseImagePrompt simulates image-to-prompt functionality
func (s *AIService) ReverseImagePrompt(imagePath string) (*ExplodeResult, error) {
	return &ExplodeResult{
		Atoms: []ExtractedAtom{
			{Value: "masterpiece", Label: "masterpiece", Type: "Positive", CategoryID: 1, CategoryName: "quality", IsNew: false},
			{Value: "best quality", Label: "best quality", Type: "Positive", CategoryID: 1, CategoryName: "quality", IsNew: false},
			{Value: "1girl", Label: "1girl", Type: "Positive", CategoryID: 2, CategoryName: "character", IsNew: false},
		},
		RawPrompt: "masterpiece, best quality, 1girl",
	}, nil
}

// SaveAIConfig saves AI configuration
func (s *AIService) SaveAIConfig(config *AIConfig) error {
	if config.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	return nil
}

// GetAIConfig retrieves AI configuration
func (s *AIService) GetAIConfig() (*AIConfig, error) {
	return &AIConfig{
		Provider: "ollama",
		Endpoint: "http://localhost:11434/v1",
		Model:    "llama2",
	}, nil
}
