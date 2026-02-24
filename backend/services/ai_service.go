package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"promptmaster/backend/logger"
	"promptmaster/backend/models"

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
	Provider     string `json:"provider"`
	ProviderType string `json:"provider_type"` // ollama, openai-compatible, etc.
	APIKey       string `json:"api_key"`
	Endpoint     string `json:"endpoint"`
	Model        string `json:"model"`
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
	CategoryID   uint     `json:"category"`                // 分类ID
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

// ========== 默认 Prompt 模板 ==========

// DefaultExplodeSystemPrompt 拆解提示词默认系统模板
const DefaultExplodeSystemPrompt = `你是一个专业的 AI 图像生成提示词拆解专家。请将用户提供的提示词拆解为原子词列表。

可用分类及其 ID：
{{category_list}}

规则：
1. 每个组件应该是一个单一的概念或属性
2. value 字段使用英文（小写，多词用下划线连接）
3. label 字段提供中文翻译
4. type 字段分类为 Positive 或 Negative
5. category 字段根据可用分类选择最合适的 ID（数字），可选分类：{{category_names}}
6. 包含相关的近义词
7. category 必须返回数字（ID），不要返回字符串

请严格返回 JSON 格式：
{"atoms": [{"value": "masterpiece", "label": "杰作", "type": "Positive", "category": 1, "synonyms": ["best quality"]}]}`

// DefaultExplodeUserPromptTemplate 拆解提示词默认用户模板
const DefaultExplodeUserPromptTemplate = `{{input}}`

// DefaultOptimizeSystemPrompt 优化提示词默认系统模板
const DefaultOptimizeSystemPrompt = `你是一个专业的 AI 图像生成提示词优化专家。请优化用户提供的提示词，提升其质量和表达效果。

优化原则：
1. 保持原意不变的情况下提升描述质量
2. 添加必要的高质量修饰词
3. 优化词汇顺序，重要的放前面
4. 去除重复或冗余的表达
5. 确保语法正确

请严格返回 JSON 格式：
{
  "optimized": "优化后的提示词",
  "changes": ["修改说明1", "修改说明2"],
  "suggestions": ["建议1", "建议2"]
}`

// DefaultOptimizeUserPromptTemplate 优化提示词默认用户模板
const DefaultOptimizeUserPromptTemplate = `{{input}}`

// DefaultTranslateSystemPrompt 翻译提示词默认系统模板
const DefaultTranslateSystemPrompt = `你是一个专业的 AI 图像生成提示词翻译专家。请将用户提供的中文提示词翻译成高质量的英文。

翻译原则：
1. 使用 AI 图像生成领域常用的专业术语
2. 保持提示词的顺序和结构
3. 使用英文逗号分隔
4. 确保翻译准确且地道
5. 提取关键概念作为关键词

请严格返回 JSON 格式：
{
  "translation": "英文翻译结果",
  "keywords": ["关键词1", "关键词2"],
  "notes": "翻译备注说明"
}`

// DefaultTranslateUserPromptTemplate 翻译提示词默认用户模板
const DefaultTranslateUserPromptTemplate = `{{input}}`

// DefaultAnalyzeSystemPrompt 分析提示词默认系统模板
const DefaultAnalyzeSystemPrompt = `你是一个专业的 AI 图像生成提示词分析专家。请分析用户提供的提示词的结构和效果。

分析维度：
1. subject - 主体内容描述
2. style - 艺术风格
3. quality - 质量相关词汇
4. lighting - 光照效果
5. other - 其他重要元素

请严格返回 JSON 格式：
{
  "analysis": {
    "subject": "主体分析",
    "style": "风格分析",
    "quality": "质量分析",
    "lighting": "光照分析",
    "other": "其他分析"
  },
  "issues": ["问题1", "问题2"],
  "suggestions": ["建议1", "建议2"]
}`

// DefaultAnalyzeUserPromptTemplate 分析提示词默认用户模板
const DefaultAnalyzeUserPromptTemplate = `{{input}}`

// shouldCallAI 检查是否应该调用AI API
// Ollama 不需要 API Key，其他提供商需要
func (s *AIService) shouldCallAI(config *AIConfig) bool {
	if config == nil {
		return false
	}
	// Ollama 本地服务不需要 API Key
	if config.ProviderType == "ollama" || config.Provider == "ollama" {
		return config.Endpoint != ""
	}
	// 其他提供商需要 API Key
	return config.APIKey != ""
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
func (s *AIService) ExplodePrompt(prompt string, categories []string, categoryMap map[string]uint, config *AIConfig, systemPrompt string, userPromptTemplate string) (*ExplodeResult, error) {
	if config == nil || !s.shouldCallAI(config) {
		return s.ruleBasedExplosion(prompt, categoryMap)
	}
	return s.aiBasedExplosion(prompt, categories, categoryMap, config, systemPrompt, userPromptTemplate)
}

// OptimizePrompt uses AI to optimize a prompt
func (s *AIService) OptimizePrompt(prompt string, config *AIConfig, systemPrompt string, userPromptTemplate string) (*OptimizeResult, error) {
	if config == nil || !s.shouldCallAI(config) {
		return s.ruleBasedOptimization(prompt)
	}
	return s.aiBasedOptimization(prompt, config, systemPrompt, userPromptTemplate)
}

// TranslatePrompt uses AI to translate a prompt
func (s *AIService) TranslatePrompt(prompt string, config *AIConfig, systemPrompt string, userPromptTemplate string) (*TranslateResult, error) {
	if config == nil || !s.shouldCallAI(config) {
		return s.ruleBasedTranslation(prompt)
	}
	return s.aiBasedTranslation(prompt, config, systemPrompt, userPromptTemplate)
}

// AnalyzePrompt uses AI to analyze a prompt
func (s *AIService) AnalyzePrompt(prompt string, config *AIConfig, systemPrompt string, userPromptTemplate string) (*AnalyzeResult, error) {
	if config == nil || !s.shouldCallAI(config) {
		return s.ruleBasedAnalysis(prompt)
	}
	return s.aiBasedAnalysis(prompt, config, systemPrompt, userPromptTemplate)
}

// renderTemplate 渲染模板，替换变量
func (s *AIService) renderTemplate(template string, variables map[string]string) string {
	result := template
	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
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
func (s *AIService) aiBasedExplosion(prompt string, categories []string, categoryMap map[string]uint, config *AIConfig, customSystemPrompt string, userPromptTemplate string) (*ExplodeResult, error) {
	// 如果 categoryMap 为空，使用默认分类
	if len(categoryMap) == 0 {
		return nil, errors.New("CategoryMap is empty")
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

	// 确定使用的 systemPrompt
	systemPrompt := customSystemPrompt
	if systemPrompt == "" {
		systemPrompt = DefaultExplodeSystemPrompt
	}
	systemPrompt = fmt.Sprintf(customSystemPrompt, categoryList)
	// 渲染 systemPrompt 模板变量
	systemPrompt = s.renderTemplate(systemPrompt, map[string]string{
		"category_list":  categoryList,
		"category_names": categoryNameList,
	})

	// 确定使用的 userPromptTemplate
	userTemplate := userPromptTemplate
	if userTemplate == "" {
		userTemplate = DefaultExplodeUserPromptTemplate
	}

	// 渲染用户提示词
	userPrompt := s.renderTemplate(userTemplate, map[string]string{
		"input": prompt,
	})

	response, err := s.callAIAPI(config, systemPrompt, userPrompt)
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
func (s *AIService) aiBasedOptimization(prompt string, config *AIConfig, customSystemPrompt string, userPromptTemplate string) (*OptimizeResult, error) {
	// 确定使用的 systemPrompt
	systemPrompt := customSystemPrompt
	if systemPrompt == "" {
		systemPrompt = DefaultOptimizeSystemPrompt
	}

	// 确定使用的 userPromptTemplate
	userTemplate := userPromptTemplate
	if userTemplate == "" {
		userTemplate = DefaultOptimizeUserPromptTemplate
	}

	// 渲染用户提示词
	userPrompt := s.renderTemplate(userTemplate, map[string]string{
		"input": prompt,
	})

	response, err := s.callAIAPI(config, systemPrompt, userPrompt)
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
func (s *AIService) aiBasedTranslation(prompt string, config *AIConfig, customSystemPrompt string, userPromptTemplate string) (*TranslateResult, error) {
	// 确定使用的 systemPrompt
	systemPrompt := customSystemPrompt
	if systemPrompt == "" {
		systemPrompt = DefaultTranslateSystemPrompt
	}

	// 确定使用的 userPromptTemplate
	userTemplate := userPromptTemplate
	if userTemplate == "" {
		userTemplate = DefaultTranslateUserPromptTemplate
	}

	// 渲染用户提示词
	userPrompt := s.renderTemplate(userTemplate, map[string]string{
		"input": prompt,
	})

	response, err := s.callAIAPI(config, systemPrompt, userPrompt)
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
func (s *AIService) aiBasedAnalysis(prompt string, config *AIConfig, customSystemPrompt string, userPromptTemplate string) (*AnalyzeResult, error) {
	// 确定使用的 systemPrompt
	systemPrompt := customSystemPrompt
	if systemPrompt == "" {
		systemPrompt = DefaultAnalyzeSystemPrompt
	}

	// 确定使用的 userPromptTemplate
	userTemplate := userPromptTemplate
	if userTemplate == "" {
		userTemplate = DefaultAnalyzeUserPromptTemplate
	}

	// 渲染用户提示词
	userPrompt := s.renderTemplate(userTemplate, map[string]string{
		"input": prompt,
	})

	response, err := s.callAIAPI(config, systemPrompt, userPrompt)
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
	// Ollama 使用不同的 API 格式
	if config.ProviderType == "ollama" || config.Provider == "ollama" {
		return s.callOllamaAPI(config, systemPrompt, userPrompt)
	}
	return s.callOpenAICompatibleAPI(config, systemPrompt, userPrompt)
}

// callOllamaAPI calls Ollama API
func (s *AIService) callOllamaAPI(config *AIConfig, systemPrompt, userPrompt string) (string, error) {
	startTime := time.Now()

	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:11434"
	}

	model := config.Model
	if model == "" {
		model = "llama2"
	}

	// Ollama 的 /api/generate 格式
	requestBody := map[string]interface{}{
		"model":  model,
		"prompt": systemPrompt + "\n\n" + userPrompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.3,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 记录请求日志
	if s.logger != nil {
		s.logger.LogAIRequest(config.Provider, endpoint+"/api/generate", model, requestBody)
	}

	req, err := http.NewRequest("POST", endpoint+"/api/generate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)

	// 记录响应日志
	duration := time.Since(startTime)
	if s.logger != nil {
		s.logger.LogAIResponse(resp.StatusCode, bodyString, duration)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API error (%d): %s", resp.StatusCode, bodyString)
	}

	// 解析 Ollama 响应
	var ollamaResponse struct {
		Response string `json:"response"`
	}

	if err := json.Unmarshal(bodyBytes, &ollamaResponse); err != nil {
		return "", err
	}

	return ollamaResponse.Response, nil
}

// callOpenAICompatibleAPI calls OpenAI compatible API
func (s *AIService) callOpenAICompatibleAPI(config *AIConfig, systemPrompt, userPrompt string) (string, error) {
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

// SaveAIConfig saves AI provider configuration to database
func (s *AIService) SaveAIConfig(config *AIConfig) error {
	if config.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	var provider models.AIProviderConfig
	result := s.db.Where("provider = ?", config.Provider).First(&provider)

	if result.Error != nil {
		// Create new provider
		provider = models.AIProviderConfig{
			Provider: config.Provider,
			Type:     config.ProviderType,
			BaseURL:  config.Endpoint,
			APIKey:   config.APIKey,
			Model:    config.Model,
		}
		return s.db.Create(&provider).Error
	}

	// Update existing provider
	updates := map[string]interface{}{
		"type":     config.ProviderType,
		"base_url": config.Endpoint,
		"api_key":  config.APIKey,
		"model":    config.Model,
	}
	return s.db.Model(&provider).Updates(updates).Error
}

// GetAIConfig retrieves current AI configuration from database
func (s *AIService) GetAIConfig() (*AIConfig, error) {
	// Get current settings
	var settings models.AISettings
	result := s.db.First(&settings)
	if result.Error != nil {
		// Return default if no settings
		return &AIConfig{
			Provider: "ollama",
			Endpoint: "http://localhost:11434/v1",
			Model:    "llama2",
		}, nil
	}

	// Get current provider config
	var provider models.AIProviderConfig
	result = s.db.Where("provider = ?", settings.CurrentProvider).First(&provider)
	if result.Error != nil {
		// Return default if provider not found
		return &AIConfig{
			Provider: "ollama",
			Endpoint: "http://localhost:11434/v1",
			Model:    "llama2",
		}, nil
	}

	return &AIConfig{
		Provider:     provider.Provider,
		ProviderType: provider.Type,
		APIKey:       provider.APIKey,
		Endpoint:     provider.BaseURL,
		Model:        provider.Model,
	}, nil
}

// GetAllProviders retrieves all AI provider configurations
func (s *AIService) GetAllProviders() ([]models.AIProviderConfig, error) {
	var providers []models.AIProviderConfig
	result := s.db.Order("sort_order ASC").Find(&providers)
	return providers, result.Error
}

// SaveProvider saves or updates an AI provider configuration
func (s *AIService) SaveProvider(provider *models.AIProviderConfig) error {
	if provider.Provider == "" {
		return fmt.Errorf("provider id is required")
	}

	var existing models.AIProviderConfig
	result := s.db.Where("provider = ?", provider.Provider).First(&existing)

	if result.Error != nil {
		// Create new
		return s.db.Create(provider).Error
	}

	// Update
	return s.db.Model(&existing).Updates(map[string]interface{}{
		"name":       provider.Name,
		"type":       provider.Type,
		"base_url":   provider.BaseURL,
		"api_key":    provider.APIKey,
		"model":      provider.Model,
		"models":     provider.Models,
		"headers":    provider.Headers,
		"enabled":    provider.Enabled,
		"is_custom":  provider.IsCustom,
		"sort_order": provider.SortOrder,
	}).Error
}

// DeleteProvider deletes an AI provider configuration
func (s *AIService) DeleteProvider(providerID string) error {
	return s.db.Where("provider = ? AND is_default = ?", providerID, false).Delete(&models.AIProviderConfig{}).Error
}

// SetCurrentProvider sets the current active provider
func (s *AIService) SetCurrentProvider(providerID string) error {
	var settings models.AISettings
	result := s.db.First(&settings)

	if result.Error != nil {
		// Create new settings
		return s.db.Create(&models.AISettings{
			CurrentProvider: providerID,
		}).Error
	}

	// Update
	return s.db.Model(&settings).Update("current_provider", providerID).Error
}

// GetAllPromptTemplates retrieves all AI prompt templates
func (s *AIService) GetAllPromptTemplates() ([]models.AIPromptTemplate, error) {
	var templates []models.AIPromptTemplate
	result := s.db.Order("template_id ASC").Find(&templates)
	return templates, result.Error
}

// GetPromptTemplate retrieves a prompt template by ID
func (s *AIService) GetPromptTemplate(templateID string) (*models.AIPromptTemplate, error) {
	var template models.AIPromptTemplate
	result := s.db.Where("template_id = ?", templateID).First(&template)
	if result.Error != nil {
		return nil, result.Error
	}
	return &template, nil
}

// SavePromptTemplate saves or updates a prompt template
func (s *AIService) SavePromptTemplate(template *models.AIPromptTemplate) error {
	if template.TemplateID == "" {
		return fmt.Errorf("template id is required")
	}

	var existing models.AIPromptTemplate
	result := s.db.Where("template_id = ?", template.TemplateID).First(&existing)

	if result.Error != nil {
		// Create new
		return s.db.Create(template).Error
	}

	// Update
	return s.db.Model(&existing).Updates(map[string]interface{}{
		"name":                 template.Name,
		"description":          template.Description,
		"system_prompt":        template.SystemPrompt,
		"user_prompt_template": template.UserPromptTemplate,
		"temperature":          template.Temperature,
		"response_format":      template.ResponseFormat,
		"is_custom":            template.IsCustom,
	}).Error
}

// DeletePromptTemplate deletes a prompt template
func (s *AIService) DeletePromptTemplate(templateID string) error {
	return s.db.Where("template_id = ? AND is_default = ?", templateID, false).Delete(&models.AIPromptTemplate{}).Error
}

// SetCurrentPrompt sets the current active prompt template
func (s *AIService) SetCurrentPrompt(templateID string) error {
	var settings models.AISettings
	result := s.db.First(&settings)

	if result.Error != nil {
		// Create new settings
		return s.db.Create(&models.AISettings{
			CurrentPrompt: templateID,
		}).Error
	}

	// Update
	return s.db.Model(&settings).Update("current_prompt", templateID).Error
}

// GetCurrentSettings retrieves current AI settings
func (s *AIService) GetCurrentSettings() (*models.AISettings, error) {
	var settings models.AISettings
	result := s.db.First(&settings)
	if result.Error != nil {
		// Return default
		return &models.AISettings{
			CurrentProvider: "ollama",
			CurrentPrompt:   "explode",
		}, nil
	}
	return &settings, nil
}
