package models

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// JSON type for storing JSON in SQLite
type JSON map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// StringSlice type for storing string arrays
type StringSlice []string

// Value implements the driver.Valuer interface
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

// Scan implements the sql.Scanner interface
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return nil
	}
	
	str = strings.TrimSpace(str)
	
	// Handle empty string
	if str == "" {
		*s = StringSlice{}
		return nil
	}
	
	// Try to parse as JSON array
	if strings.HasPrefix(str, "[") {
		return json.Unmarshal([]byte(str), s)
	}
	
	// Handle legacy single string format - treat as single-element array
	// Also handle quoted strings like '"value"'
	str = strings.Trim(str, `"`)
	if str != "" {
		*s = StringSlice{str}
	} else {
		*s = StringSlice{}
	}
	return nil
}

// Category represents a hierarchical category
type Category struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	ParentID   uint      `gorm:"default:0;index" json:"parent_id"`
	Type       string    `gorm:"size:20;not null" json:"type"` // ATOM or PRESET
	SortOrder  int       `gorm:"default:0" json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Atom represents an atomic prompt word
type Atom struct {
	ID          uint        `gorm:"primarykey" json:"id"`
	Value       string      `gorm:"size:200;not null;uniqueIndex" json:"value"` // English word
	Label       string      `gorm:"size:200" json:"label"`                      // Chinese label
	Synonyms    StringSlice `gorm:"type:text" json:"synonyms"`                  // JSON array of synonyms
	Type        string      `gorm:"size:20;not null;default:'Positive'" json:"type"` // Positive or Negative
	CategoryID  uint        `gorm:"index" json:"category_id"`
	UsageCount  int         `gorm:"default:0" json:"usage_count"`
	LastUsedAt  *time.Time  `json:"last_used_at"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	
	// Relationships
	Category    Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// Preset represents a prompt preset
type Preset struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	Title           string    `gorm:"size:200;not null" json:"title"`
	CategoryID      uint      `gorm:"default:0;index" json:"category_id"`
	CurrentVersion  int       `gorm:"default:0" json:"current_version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	IsDeleted       bool      `gorm:"default:false" json:"is_deleted"`
	
	// Relationships
	Category        Category        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Versions        []PresetVersion `gorm:"foreignKey:PresetID" json:"versions,omitempty"`
}

// PresetVersion represents a version snapshot of a preset
type PresetVersion struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	PresetID       uint      `gorm:"not null;index" json:"preset_id"`
	VersionNum     int       `gorm:"not null" json:"version_num"` // V1=1, V2=2, etc.
	Snapshot       JSON      `gorm:"type:text" json:"snapshot"`   // Full JSON snapshot
	ThumbnailPath  string    `gorm:"size:500" json:"thumbnail_path"`
	IsStarred      bool      `gorm:"default:false" json:"is_starred"`
	CreatedAt      time.Time `json:"created_at"`
	DiffStats      string    `gorm:"size:20" json:"diff_stats"`   // e.g., "+2/-1"
	
	// Relationships
	Preset         Preset    `gorm:"foreignKey:PresetID" json:"preset,omitempty"`
	Previews       []Preview `gorm:"foreignKey:VersionID" json:"previews,omitempty"`
}

// SnapshotData represents the structure of the JSON snapshot
type SnapshotData struct {
	PosText      string                 `json:"pos_text"`
	NegText      string                 `json:"neg_text"`
	Params       map[string]interface{} `json:"params"`
	PreviewPaths []string               `json:"preview_paths"`
	AtomIDs      []uint                 `json:"atom_ids"`
}

// Preview represents a preview image
type Preview struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	PresetID   uint      `gorm:"index" json:"preset_id"`
	VersionID  *uint     `gorm:"index" json:"version_id,omitempty"`
	FilePath   string    `gorm:"size:500;not null" json:"file_path"`
	CreatedAt  time.Time `json:"created_at"`
}

// UsageStat tracks usage statistics for atoms
type UsageStat struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	AtomID     uint      `gorm:"uniqueIndex;not null" json:"atom_id"`
	UseCount   int       `gorm:"default:0" json:"use_count"`
	LastUsedAt time.Time `json:"last_used_at"`
}

// ToSnapshotData converts JSON snapshot to SnapshotData
func (pv *PresetVersion) ToSnapshotData() (*SnapshotData, error) {
	if pv.Snapshot == nil {
		return &SnapshotData{}, nil
	}
	
	data := &SnapshotData{}
	if posText, ok := pv.Snapshot["pos_text"].(string); ok {
		data.PosText = posText
	}
	if negText, ok := pv.Snapshot["neg_text"].(string); ok {
		data.NegText = negText
	}
	if params, ok := pv.Snapshot["params"].(map[string]interface{}); ok {
		data.Params = params
	}
	if paths, ok := pv.Snapshot["preview_paths"].([]interface{}); ok {
		for _, p := range paths {
			if str, ok := p.(string); ok {
				data.PreviewPaths = append(data.PreviewPaths, str)
			}
		}
	}
	if ids, ok := pv.Snapshot["atom_ids"].([]interface{}); ok {
		for _, id := range ids {
			if num, ok := id.(float64); ok {
				data.AtomIDs = append(data.AtomIDs, uint(num))
			}
		}
	}
	return data, nil
}

// SetSnapshotData sets the snapshot from SnapshotData
func (pv *PresetVersion) SetSnapshotData(data *SnapshotData) error {
	pv.Snapshot = JSON{
		"pos_text":      data.PosText,
		"neg_text":      data.NegText,
		"params":        data.Params,
		"preview_paths": data.PreviewPaths,
		"atom_ids":      data.AtomIDs,
	}
	return nil
}

// AIProviderConfig represents AI provider configuration stored in database
type AIProviderConfig struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Provider  string    `gorm:"size:50;not null;uniqueIndex" json:"provider"` // 提供商ID，如 openai, ollama
	Name      string    `gorm:"size:100" json:"name"`                        // 显示名称
	Type      string    `gorm:"size:50" json:"type"`                          // 类型: openai-compatible, ollama
	BaseURL   string    `gorm:"size:500" json:"base_url"`                     // API地址
	APIKey    string    `gorm:"size:500" json:"api_key"`                      // API密钥（加密存储）
	Model     string    `gorm:"size:100" json:"model"`                        // 默认模型
	Models    StringSlice `gorm:"type:text" json:"models"`                    // 可用模型列表
	Headers   JSON      `gorm:"type:text" json:"headers"`                     // 自定义请求头
	Enabled   bool      `gorm:"default:false" json:"enabled"`                 // 是否启用
	IsCustom  bool      `gorm:"default:false" json:"is_custom"`               // 是否用户自定义
	SortOrder int       `gorm:"default:0" json:"sort_order"`                  // 排序
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AIPromptTemplate represents AI prompt template stored in database
type AIPromptTemplate struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	TemplateID       string    `gorm:"size:50;not null;uniqueIndex" json:"template_id"` // 模板ID，如 explode, optimize
	Name             string    `gorm:"size:100" json:"name"`                           // 显示名称
	Description      string    `gorm:"size:500" json:"description"`                    // 描述
	SystemPrompt     string    `gorm:"type:text" json:"system_prompt"`                 // 系统提示词
	UserPromptTemplate string `gorm:"type:text" json:"user_prompt_template"`          // 用户提示词模板
	Temperature      float64   `gorm:"default:0.7" json:"temperature"`                 // 温度参数
	ResponseFormat   string    `gorm:"size:20;default:'json'" json:"response_format"`  // 响应格式
	IsCustom         bool      `gorm:"default:false" json:"is_custom"`                 // 是否用户自定义
	IsDefault        bool      `gorm:"default:false" json:"is_default"`                // 是否系统默认（不可删除）
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// AISettings represents global AI settings
type AISettings struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	CurrentProvider  string    `gorm:"size:50" json:"current_provider"`  // 当前选中的提供商ID
	CurrentPrompt    string    `gorm:"size:50" json:"current_prompt"`    // 当前选中的提示词ID
	UpdatedAt        time.Time `json:"updated_at"`
}
