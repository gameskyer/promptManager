package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

// TranslationDirection specifies the direction of translation.
type TranslationDirection string

const (
	DirCNtoEN TranslationDirection = "cn2en" // Chinese → English
	DirENtoCN TranslationDirection = "en2cn" // English → Chinese
)

type PromptWord struct {
	ID           string `json:"id"`
	Original     string `json:"original"`
	Translated   string `json:"translated"`
	NeedsTrans   bool   `json:"needs_trans"`
	IsTranslated bool   `json:"is_translated"`
}

type TransSessionService struct {
	aiService *AIService
}

func NewTransSessionService(aiService *AIService) *TransSessionService {
	return &TransSessionService{aiService: aiService}
}

// SplitPrompt splits raw prompt text by both Chinese (，；) and English (,) commas.
// direction determines which words need translation.
func (s *TransSessionService) SplitPrompt(rawText string, direction TranslationDirection) []PromptWord {
	delimiterRegex := regexp.MustCompile(`[,，;；]`)
	parts := delimiterRegex.Split(rawText, -1)

	var words []PromptWord
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		hasChinese := s.containsChinese(trimmed)
		needsTrans := s.needsTranslation(trimmed, direction, hasChinese)

		words = append(words, PromptWord{
			ID:           uuid.New().String(),
			Original:     trimmed,
			Translated:   trimmed,
			NeedsTrans:   needsTrans,
			IsTranslated: !needsTrans,
		})
	}
	return words
}

// needsTranslation determines if a word needs translation based on direction.
func (s *TransSessionService) needsTranslation(text string, direction TranslationDirection, hasChinese bool) bool {
	switch direction {
	case DirCNtoEN:
		return hasChinese
	case DirENtoCN:
		return !hasChinese && s.containsLatin(text)
	default:
		return hasChinese
	}
}

// TranslateWords translates all words that need translation but aren't yet translated.
func (s *TransSessionService) TranslateWords(words []PromptWord, direction TranslationDirection) ([]PromptWord, error) {
	toTranslate := make([]PromptWord, 0)
	translateIdx := make(map[string]int)
	for i, w := range words {
		if w.NeedsTrans && !w.IsTranslated {
			toTranslate = append(toTranslate, w)
			translateIdx[w.Original] = i
		}
	}

	if len(toTranslate) == 0 {
		return words, nil
	}

	translations, err := s.batchTranslate(toTranslate, direction)
	if err != nil {
		return words, err
	}

	result := make([]PromptWord, len(words))
	copy(result, words)
	for _, tw := range toTranslate {
		if translated, ok := translations[tw.Original]; ok {
			idx := translateIdx[tw.Original]
			result[idx].Translated = translated
			result[idx].IsTranslated = true
		}
	}

	return result, nil
}

// TranslateSingle translates a single word.
func (s *TransSessionService) TranslateSingle(text string, direction TranslationDirection) (string, error) {
	needsTrans := s.needsTranslation(text, direction, s.containsChinese(text))
	words := []PromptWord{
		{Original: text, NeedsTrans: needsTrans, IsTranslated: !needsTrans},
	}
	result, err := s.TranslateWords(words, direction)
	if err != nil {
		return text, err
	}
	return result[0].Translated, nil
}

func (s *TransSessionService) batchTranslate(toTranslate []PromptWord, direction TranslationDirection) (map[string]string, error) {
	if len(toTranslate) == 0 {
		return map[string]string{}, nil
	}

	if s.aiService == nil {
		return map[string]string{}, fmt.Errorf("AI service not available")
	}

	var texts []string
	for _, w := range toTranslate {
		texts = append(texts, w.Original)
	}
	wordsJSON, _ := json.Marshal(texts)

	var systemPrompt string
	switch direction {
	case DirCNtoEN:
		systemPrompt = s.cn2enSystemPrompt()
	case DirENtoCN:
		systemPrompt = s.en2cnSystemPrompt()
	default:
		systemPrompt = s.cn2enSystemPrompt()
	}

	config, err := s.aiService.GetAIConfig()
	if err != nil {
		config = &AIConfig{}
	}

	if !s.aiService.shouldCallAI(config) {
		result := make(map[string]string)
		for _, w := range toTranslate {
			result[w.Original] = w.Original
		}
		return result, nil
	}

	response, err := s.aiService.callAIAPI(config, systemPrompt, string(wordsJSON))
	if err != nil {
		return nil, fmt.Errorf("翻译请求失败: %w", err)
	}

	translations, err := s.parseTranslationResponse(response)
	if err != nil {
		return nil, fmt.Errorf("解析翻译结果失败: %w", err)
	}

	return translations, nil
}

func (s *TransSessionService) cn2enSystemPrompt() string {
	return `你是一个专注于翻译AI绘画提示词标签的专业翻译机器。

## 任务
将中文提示词标签翻译为英文。每个标签是一个独立的视觉概念词。

## 翻译规则
1. 使用AI绘画社区（Stable Diffusion/NovelAI/Danbooru）的标准英文标签
2. 多词使用下划线连接（如：蓝色天空 -> blue_sky）
3. 保持标签的准确性和简洁性
4. 专有名词保留原样（如角色名、作品名）
5. 一个词或女孩等计数使用数字+英文（如：一个女孩 -> 1girl，两个男孩 -> 2boys）

## 输出格式
严格返回以下JSON格式，不要包含任何其他文字：
{
  "translations": {
    "中文词1": "english_tag_1",
    "中文词2": "english_tag_2"
  }
}

请确保：
- 每个输入的中文词都有对应的翻译
- JSON格式完全合法
- 不要输出任何Markdown标记`
}

func (s *TransSessionService) en2cnSystemPrompt() string {
	return `你是一个专注于翻译AI绘画提示词标签的专业翻译机器。

## 任务
将英文提示词标签翻译为中文。每个标签是一个独立的视觉概念词。

## 翻译规则
1. 使用AI绘画社区常用的中文表达
2. 翻译为简洁准确的中文词汇或短语
3. 专有名词（如角色名、作品名）保留原样或使用公认译名
4. 风格、技法类词汇使用绘画领域通用中文术语（如：masterpiece -> 杰作，oil painting -> 油画）
5. 质量词给出简洁中文注释

## 输出格式
严格返回以下JSON格式，不要包含任何其他文字：
{
  "translations": {
    "masterpiece": "杰作",
    "best_quality": "最佳质量",
    "1girl": "一个女孩"
  }
}

请确保：
- 每个输入的英文词都有对应的翻译
- JSON格式完全合法
- 不要输出任何Markdown标记`
}

// parseTranslationResponse extracts translation mapping from AI response.
func (s *TransSessionService) parseTranslationResponse(response string) (map[string]string, error) {
	jsonStr := extractJSONStr(response)

	var result struct {
		Translations map[string]string `json:"translations"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		var plainMap map[string]string
		if err2 := json.Unmarshal([]byte(jsonStr), &plainMap); err2 != nil {
			return nil, fmt.Errorf("无法解析翻译响应: %s", jsonStr[:min(len(jsonStr), 100)])
		}
		return plainMap, nil
	}

	return result.Translations, nil
}

func (s *TransSessionService) containsChinese(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func (s *TransSessionService) containsLatin(text string) bool {
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}

func extractJSONStr(content string) string {
	re := regexp.MustCompile("`{3}(?:json)?\\s*([\\s\\S]*?)`{3}")
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start != -1 && end != -1 && end > start {
		return content[start : end+1]
	}

	return strings.TrimSpace(content)
}
