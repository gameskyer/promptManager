package models

import (
	"fmt"
	"promptmaster/backend/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection and migrates the schema
func InitDB() (*gorm.DB, error) {
	dbPath := config.GetDBPath()
	
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Auto migrate all models
	err = db.AutoMigrate(
		&Category{},
		&Atom{},
		&Preset{},
		&PresetVersion{},
		&Preview{},
		&UsageStat{},
		&AIProviderConfig{},
		&AIPromptTemplate{},
		&AISettings{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	
	// Initialize default categories
	initDefaultCategories(db)
	
	// Initialize default AI configs
	initDefaultAIConfigs(db)
	
	return db, nil
}

// initDefaultCategories creates default categories if none exist
func initDefaultCategories(db *gorm.DB) {
	var count int64
	db.Model(&Category{}).Count(&count)
	if count > 0 {
		return
	}

	// Create default categories
	categories := []Category{
		// 一级分类 - 提示词库
		{Name: "人物", ParentID: 0, Type: "ATOM", SortOrder: 1},
		{Name: "场景", ParentID: 0, Type: "ATOM", SortOrder: 2},
		{Name: "风格", ParentID: 0, Type: "ATOM", SortOrder: 3},
		{Name: "质量", ParentID: 0, Type: "ATOM", SortOrder: 4},
		
		// 一级分类 - 预设库（只保留一个）
		{Name: "预设库", ParentID: 0, Type: "PRESET", SortOrder: 10},
	}

	for i := range categories {
		db.Create(&categories[i])
	}

	// 创建二级分类示例（人物下的分类）
	var personCategory Category
	db.Where("name = ? AND type = ?", "人物", "ATOM").First(&personCategory)
	
	if personCategory.ID > 0 {
		subCategories := []Category{
			{Name: "发型", ParentID: personCategory.ID, Type: "ATOM", SortOrder: 1},
			{Name: "眼睛", ParentID: personCategory.ID, Type: "ATOM", SortOrder: 2},
			{Name: "服装", ParentID: personCategory.ID, Type: "ATOM", SortOrder: 3},
			{Name: "姿势", ParentID: personCategory.ID, Type: "ATOM", SortOrder: 4},
			{Name: "表情", ParentID: personCategory.ID, Type: "ATOM", SortOrder: 5},
		}
		for i := range subCategories {
			db.Create(&subCategories[i])
		}
	}
}

// initDefaultAIConfigs creates default AI provider and prompt configs if none exist
func initDefaultAIConfigs(db *gorm.DB) {
	// Check if providers exist
	var providerCount int64
	db.Model(&AIProviderConfig{}).Count(&providerCount)
	
	if providerCount == 0 {
		// Create default AI providers
		providers := []AIProviderConfig{
			{
				Provider:  "openai",
				Name:      "OpenAI",
				Type:      "openai-compatible",
				BaseURL:   "https://api.openai.com/v1",
				Model:     "gpt-3.5-turbo",
				Models:    StringSlice{"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo"},
				Enabled:   false,
				IsCustom:  false,
				SortOrder: 1,
			},
			{
				Provider:  "deepseek",
				Name:      "DeepSeek",
				Type:      "openai-compatible",
				BaseURL:   "https://api.deepseek.com/v1",
				Model:     "deepseek-chat",
				Models:    StringSlice{"deepseek-chat", "deepseek-coder"},
				Enabled:   false,
				IsCustom:  false,
				SortOrder: 2,
			},
			{
				Provider:  "kimi",
				Name:      "Kimi (Moonshot)",
				Type:      "openai-compatible",
				BaseURL:   "https://api.moonshot.cn/v1",
				Model:     "moonshot-v1-8k",
				Models:    StringSlice{"moonshot-v1-8k", "moonshot-v1-32k", "moonshot-v1-128k"},
				Enabled:   false,
				IsCustom:  false,
				SortOrder: 3,
			},
			{
				Provider:  "ollama",
				Name:      "Ollama (本地)",
				Type:      "ollama",
				BaseURL:   "http://localhost:11434",
				Model:     "llama2",
				Models:    StringSlice{"llama2", "mistral", "codellama", "vicuna"},
				Enabled:   true,
				IsCustom:  false,
				SortOrder: 4,
			},
		}
		
		for i := range providers {
			db.Create(&providers[i])
		}
	}
	
	// Check if prompt templates exist
	var promptCount int64
	db.Model(&AIPromptTemplate{}).Count(&promptCount)
	
	if promptCount == 0 {
		// Create default prompt templates
		prompts := []AIPromptTemplate{
			{
				TemplateID:       "explode",
				Name:             "拆解提示词",
				Description:      "将长段提示词拆解为原子词列表",
				SystemPrompt:     getDefaultExplodePrompt(),
				UserPromptTemplate: "{{input}}",
				Temperature:      0.3,
				ResponseFormat:   "json",
				IsCustom:         false,
				IsDefault:        true,
			},
			{
				TemplateID:       "optimize",
				Name:             "优化提示词",
				Description:      "优化提示词的质量和表达",
				SystemPrompt:     getDefaultOptimizePrompt(),
				UserPromptTemplate: "{{input}}",
				Temperature:      0.4,
				ResponseFormat:   "json",
				IsCustom:         false,
				IsDefault:        true,
			},
			{
				TemplateID:       "translate",
				Name:             "翻译提示词",
				Description:      "将中文提示词翻译为英文",
				SystemPrompt:     getDefaultTranslatePrompt(),
				UserPromptTemplate: "{{input}}",
				Temperature:      0.5,
				ResponseFormat:   "json",
				IsCustom:         false,
				IsDefault:        true,
			},
			{
				TemplateID:       "analyze",
				Name:             "分析提示词",
				Description:      "分析提示词的结构和效果",
				SystemPrompt:     getDefaultAnalyzePrompt(),
				UserPromptTemplate: "{{input}}",
				Temperature:      0.6,
				ResponseFormat:   "json",
				IsCustom:         false,
				IsDefault:        true,
			},
		}
		
		for i := range prompts {
			db.Create(&prompts[i])
		}
	}
	
	// Check if settings exist
	var settingsCount int64
	db.Model(&AISettings{}).Count(&settingsCount)
	
	if settingsCount == 0 {
		// Create default settings
		db.Create(&AISettings{
			CurrentProvider: "ollama",
			CurrentPrompt:   "explode",
		})
	}
}

// Default prompt content functions
func getDefaultExplodePrompt() string {
	return `你是一位专业的 AI 图像生成提示词拆解专家，擅长将复杂的提示词分解为结构化的原子词汇。

## 任务描述
将用户提供的提示词拆解为最小的语义单元（原子词），每个原子词代表一个独立的视觉概念或属性。

## 可用分类及ID
{{category_list}}

## 拆解规则

### 1. Value 字段（英文标识）
- 使用英文小写
- 多词使用下划线连接（如：blue_sky）
- 使用 Stable Diffusion 社区标准标签

### 2. Label 字段（中文说明）
- 提供简洁准确的中文翻译

### 3. Type 字段（类型标记）
- Positive: 正向提示词（期望出现的元素）
- Negative: 负向提示词（期望避免的元素）

### 4. Category 字段（分类ID）
- 必须是数字（ID），从可用分类中选择最合适的

### 5. Synonyms 字段（近义词列表）
- 提供 1-3 个常用同义词或变体

## 输出格式
{"atoms": [{"value": "masterpiece", "label": "杰作", "type": "Positive", "category": 1, "synonyms": ["best_quality"]}]}`
}

func getDefaultOptimizePrompt() string {
	return `你是一位资深的 AI 图像生成提示词工程师，专注于优化提示词以获得最佳的图像生成效果。

## 优化原则
1. 添加必要的基础质量词
2. 按重要性排序：质量词 → 主体 → 细节 → 风格 → 光照 → 背景
3. 去除语义重复的词汇
4. 使用英文逗号分隔
5. 确保风格与主体匹配

## 输出格式
{
  "optimized": "优化后的提示词",
  "changes": ["修改说明1", "修改说明2"],
  "suggestions": ["建议1", "建议2"]
}`
}

func getDefaultTranslatePrompt() string {
	return `你是一位专业的 AI 图像生成提示词翻译专家，精通中文到英文的提示词翻译。

## 翻译原则
1. 使用 AI 绘画社区广泛认可的标准标签
2. 使用英文逗号分隔
3. 标签使用小写（专有名词除外）
4. 保持原意的完整传达
5. 质量词前置

## 输出格式
{
  "translation": "英文翻译结果",
  "keywords": ["关键词1", "关键词2"],
  "notes": "翻译说明"
}`
}

func getDefaultAnalyzePrompt() string {
	return `你是一位资深的 AI 图像生成提示词分析师，擅长深度分析提示词的构成、潜在问题和改进空间。

## 分析维度
1. Subject（主体分析）
2. Style（艺术风格分析）
3. Quality（质量相关分析）
4. Lighting（光照效果分析）
5. Other（其他要素分析）

## 输出格式
{
  "analysis": {
    "subject": "主体分析...",
    "style": "风格分析...",
    "quality": "质量分析...",
    "lighting": "光照分析...",
    "other": "其他分析..."
  },
  "issues": ["问题1", "问题2"],
  "suggestions": ["建议1", "建议2"]
}`
}
