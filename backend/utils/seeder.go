package utils

import (
	"encoding/json"
	"fmt"
	"promptmaster/backend/models"
	"promptmaster/backend/services"

	"gorm.io/gorm"
)

// SeedData represents all seed data
type SeedData struct {
	Categories []CategorySeed `json:"categories"`
	Atoms      []AtomSeed     `json:"atoms"`
	Presets    []PresetSeed   `json:"presets"`
}

// CategorySeed represents a category to seed
type CategorySeed struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	ParentID uint   `json:"parent_id"`
}

// AtomSeed represents an atom to seed
type AtomSeed struct {
	Value      string   `json:"value"`
	Label      string   `json:"label"`
	Type       string   `json:"type"`
	Category   string   `json:"category"`
	Synonyms   []string `json:"synonyms"`
	UsageCount int      `json:"usage_count"`
}

// PresetSeed represents a preset to seed
type PresetSeed struct {
	Title       string            `json:"title"`
	PosText     string            `json:"pos_text"`
	NegText     string            `json:"neg_text"`
	Params      map[string]interface{} `json:"params"`
	Versions    []VersionSeed     `json:"versions"`
}

// VersionSeed represents a version to seed
type VersionSeed struct {
	VersionNum    int               `json:"version_num"`
	PosText       string            `json:"pos_text"`
	NegText       string            `json:"neg_text"`
	Params        map[string]interface{} `json:"params"`
	IsStarred     bool              `json:"is_starred"`
	DiffStats     string            `json:"diff_stats"`
}

// Seeder handles database seeding
type Seeder struct {
	db              *gorm.DB
	categoryService *services.CategoryService
	atomService     *services.AtomService
	presetService   *services.PresetService
}

// NewSeeder creates a new seeder
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		db:              db,
		categoryService: services.NewCategoryService(db),
		atomService:     services.NewAtomService(db),
		presetService:   services.NewPresetService(db),
	}
}

// SeedAll seeds all default data
func (s *Seeder) SeedAll() error {
	fmt.Println("Starting database seeding...")

	// Check if data already exists
	var categoryCount int64
	s.db.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount > 4 { // More than default categories
		fmt.Println("Database already has data, skipping seed...")
		return nil
	}

	// Seed categories
	if err := s.seedCategories(); err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}

	// Seed atoms
	if err := s.seedAtoms(); err != nil {
		return fmt.Errorf("failed to seed atoms: %w", err)
	}

	// Seed presets
	if err := s.seedPresets(); err != nil {
		return fmt.Errorf("failed to seed presets: %w", err)
	}

	fmt.Println("Database seeding completed successfully!")
	return nil
}

// SeedFromJSON seeds data from JSON string
func (s *Seeder) SeedFromJSON(jsonData string) error {
	var data SeedData
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Seed categories first
	categoryMap := make(map[string]uint)
	
	for _, cat := range data.Categories {
		category, err := s.categoryService.CreateCategory(cat.Name, cat.Type, cat.ParentID)
		if err != nil {
			fmt.Printf("Warning: failed to create category %s: %v\n", cat.Name, err)
			continue
		}
		categoryMap[cat.Name] = category.ID
		fmt.Printf("Created category: %s (ID: %d)\n", cat.Name, category.ID)
	}

	// Seed atoms
	for _, atom := range data.Atoms {
		categoryID := categoryMap[atom.Category]
		if categoryID == 0 {
			// Try to find by name
			var cat models.Category
			if err := s.db.Where("name = ?", atom.Category).First(&cat).Error; err == nil {
				categoryID = cat.ID
			}
		}
		
		_, err := s.atomService.CreateAtom(atom.Value, atom.Label, atom.Type, categoryID, atom.Synonyms)
		if err != nil {
			fmt.Printf("Warning: failed to create atom %s: %v\n", atom.Value, err)
			continue
		}
		fmt.Printf("Created atom: %s (%s)\n", atom.Value, atom.Label)
	}

	return nil
}

// GetDefaultSeedData returns the default seed data as JSON
func (s *Seeder) GetDefaultSeedData() string {
	data := SeedData{
		Categories: []CategorySeed{
			// 一级分类 - 提示词库
			{Name: "人物", Type: "ATOM", ParentID: 0},
			{Name: "场景", Type: "ATOM", ParentID: 0},
			{Name: "风格", Type: "ATOM", ParentID: 0},
			{Name: "质量", Type: "ATOM", ParentID: 0},
			{Name: "光照", Type: "ATOM", ParentID: 0},
			{Name: "道具", Type: "ATOM", ParentID: 0},
			
			// 人物子分类
			{Name: "发型", Type: "ATOM", ParentID: 1},
			{Name: "眼睛", Type: "ATOM", ParentID: 1},
			{Name: "服装", Type: "ATOM", ParentID: 1},
			{Name: "姿势", Type: "ATOM", ParentID: 1},
			{Name: "表情", Type: "ATOM", ParentID: 1},
			
			// 场景子分类
			{Name: "室内", Type: "ATOM", ParentID: 2},
			{Name: "室外", Type: "ATOM", ParentID: 2},
			{Name: "自然", Type: "ATOM", ParentID: 2},
			{Name: "建筑", Type: "ATOM", ParentID: 2},
			
			// 风格子分类
			{Name: "艺术风格", Type: "ATOM", ParentID: 3},
			{Name: "时代风格", Type: "ATOM", ParentID: 3},
			
			// 预设库
			{Name: "预设库", Type: "PRESET", ParentID: 0},
		},
		Atoms: []AtomSeed{
			// 质量
			{Value: "masterpiece", Label: "杰作", Type: "Positive", Category: "质量", Synonyms: []string{"best quality", "top quality"}, UsageCount: 500},
			{Value: "best quality", Label: "最佳质量", Type: "Positive", Category: "质量", Synonyms: []string{"high quality", "ultra quality"}, UsageCount: 480},
			{Value: "highres", Label: "高分辨率", Type: "Positive", Category: "质量", Synonyms: []string{"high resolution", "hires"}, UsageCount: 400},
			{Value: "ultra-detailed", Label: "超详细", Type: "Positive", Category: "质量", Synonyms: []string{"highly detailed"}, UsageCount: 350},
			{Value: "8k uhd", Label: "8K超高清", Type: "Positive", Category: "质量", Synonyms: []string{"8k", "uhd"}, UsageCount: 300},
			
			// 人物
			{Value: "1girl", Label: "1个女孩", Type: "Positive", Category: "人物", Synonyms: []string{"solo girl", "single girl"}, UsageCount: 600},
			{Value: "1boy", Label: "1个男孩", Type: "Positive", Category: "人物", Synonyms: []string{"solo boy"}, UsageCount: 400},
			{Value: "solo", Label: "单人", Type: "Positive", Category: "人物", Synonyms: []string{"alone", "single"}, UsageCount: 450},
			
			// 发型
			{Value: "long hair", Label: "长发", Type: "Positive", Category: "发型", Synonyms: []string{"lengthy hair", "flowing hair"}, UsageCount: 350},
			{Value: "short hair", Label: "短发", Type: "Positive", Category: "发型", Synonyms: []string{"bob cut"}, UsageCount: 280},
			{Value: "blonde hair", Label: "金发", Type: "Positive", Category: "发型", Synonyms: []string{"golden hair", "fair hair"}, UsageCount: 250},
			{Value: "black hair", Label: "黑发", Type: "Positive", Category: "发型", Synonyms: []string{"dark hair"}, UsageCount: 400},
			{Value: "brown hair", Label: "棕发", Type: "Positive", Category: "发型", Synonyms: []string{"brunette"}, UsageCount: 220},
			{Value: "white hair", Label: "白发", Type: "Positive", Category: "发型", Synonyms: []string{"silver hair"}, UsageCount: 180},
			{Value: "red hair", Label: "红发", Type: "Positive", Category: "发型", Synonyms: []string{"ginger hair"}, UsageCount: 150},
			{Value: "blue hair", Label: "蓝发", Type: "Positive", Category: "发型", Synonyms: []string{}, UsageCount: 120},
			{Value: "ponytail", Label: "马尾辫", Type: "Positive", Category: "发型", Synonyms: []string{}, UsageCount: 200},
			{Value: "twintails", Label: "双马尾", Type: "Positive", Category: "发型", Synonyms: []string{"pigtails"}, UsageCount: 180},
			
			// 眼睛
			{Value: "blue eyes", Label: "蓝眼睛", Type: "Positive", Category: "眼睛", Synonyms: []string{"azure eyes"}, UsageCount: 300},
			{Value: "red eyes", Label: "红眼睛", Type: "Positive", Category: "眼睛", Synonyms: []string{"crimson eyes"}, UsageCount: 200},
			{Value: "green eyes", Label: "绿眼睛", Type: "Positive", Category: "眼睛", Synonyms: []string{"emerald eyes"}, UsageCount: 180},
			{Value: "brown eyes", Label: "棕眼睛", Type: "Positive", Category: "眼睛", Synonyms: []string{}, UsageCount: 220},
			{Value: "purple eyes", Label: "紫眼睛", Type: "Positive", Category: "眼睛", Synonyms: []string{}, UsageCount: 150},
			
			// 服装
			{Value: "school uniform", Label: "校服", Type: "Positive", Category: "服装", Synonyms: []string{"seifuku", "uniform"}, UsageCount: 320},
			{Value: "dress", Label: "连衣裙", Type: "Positive", Category: "服装", Synonyms: []string{"gown", "frock"}, UsageCount: 280},
			{Value: "shirt", Label: "衬衫", Type: "Positive", Category: "服装", Synonyms: []string{}, UsageCount: 200},
			{Value: "skirt", Label: "裙子", Type: "Positive", Category: "服装", Synonyms: []string{}, UsageCount: 250},
			{Value: "jacket", Label: "夹克", Type: "Positive", Category: "服装", Synonyms: []string{}, UsageCount: 150},
			{Value: "hoodie", Label: "连帽衫", Type: "Positive", Category: "服装", Synonyms: []string{}, UsageCount: 120},
			
			// 姿势
			{Value: "standing", Label: "站立", Type: "Positive", Category: "姿势", Synonyms: []string{"stand"}, UsageCount: 350},
			{Value: "sitting", Label: "坐着", Type: "Positive", Category: "姿势", Synonyms: []string{"sit"}, UsageCount: 280},
			{Value: "walking", Label: "行走", Type: "Positive", Category: "姿势", Synonyms: []string{}, UsageCount: 200},
			
			// 表情
			{Value: "smile", Label: "微笑", Type: "Positive", Category: "表情", Synonyms: []string{"grin", "happy"}, UsageCount: 400},
			{Value: "laughing", Label: "大笑", Type: "Positive", Category: "表情", Synonyms: []string{}, UsageCount: 150},
			{Value: "serious", Label: "严肃", Type: "Positive", Category: "表情", Synonyms: []string{}, UsageCount: 180},
			
			// 场景
			{Value: "outdoors", Label: "户外", Type: "Positive", Category: "室外", Synonyms: []string{"outside"}, UsageCount: 300},
			{Value: "indoors", Label: "室内", Type: "Positive", Category: "室内", Synonyms: []string{"inside"}, UsageCount: 250},
			{Value: "classroom", Label: "教室", Type: "Positive", Category: "室内", Synonyms: []string{}, UsageCount: 200},
			{Value: "bedroom", Label: "卧室", Type: "Positive", Category: "室内", Synonyms: []string{}, UsageCount: 180},
			{Value: "city", Label: "城市", Type: "Positive", Category: "建筑", Synonyms: []string{"urban"}, UsageCount: 280},
			{Value: "nature", Label: "自然", Type: "Positive", Category: "自然", Synonyms: []string{}, UsageCount: 220},
			{Value: "forest", Label: "森林", Type: "Positive", Category: "自然", Synonyms: []string{}, UsageCount: 180},
			{Value: "beach", Label: "海滩", Type: "Positive", Category: "自然", Synonyms: []string{"seaside"}, UsageCount: 160},
			
			// 风格
			{Value: "anime style", Label: "动漫风格", Type: "Positive", Category: "艺术风格", Synonyms: []string{"anime"}, UsageCount: 450},
			{Value: "realistic", Label: "写实风格", Type: "Positive", Category: "艺术风格", Synonyms: []string{"photorealistic"}, UsageCount: 350},
			{Value: "chibi", Label: "Q版", Type: "Positive", Category: "艺术风格", Synonyms: []string{"deformed"}, UsageCount: 120},
			
			// 光照
			{Value: "cinematic lighting", Label: "电影光照", Type: "Positive", Category: "光照", Synonyms: []string{}, UsageCount: 250},
			{Value: "soft lighting", Label: "柔和光照", Type: "Positive", Category: "光照", Synonyms: []string{}, UsageCount: 200},
			{Value: "sunlight", Label: "阳光", Type: "Positive", Category: "光照", Synonyms: []string{}, UsageCount: 180},
			
			// 负向提示词
			{Value: "low quality", Label: "低质量", Type: "Negative", Category: "质量", Synonyms: []string{"bad quality"}, UsageCount: 500},
			{Value: "bad anatomy", Label: "糟糕的人体结构", Type: "Negative", Category: "人物", Synonyms: []string{}, UsageCount: 400},
			{Value: "bad hands", Label: "糟糕的手", Type: "Negative", Category: "人物", Synonyms: []string{}, UsageCount: 380},
			{Value: "text", Label: "文字", Type: "Negative", Category: "质量", Synonyms: []string{}, UsageCount: 350},
			{Value: "watermark", Label: "水印", Type: "Negative", Category: "质量", Synonyms: []string{}, UsageCount: 320},
			{Value: "blurry", Label: "模糊", Type: "Negative", Category: "质量", Synonyms: []string{}, UsageCount: 300},
			{Value: "missing fingers", Label: "缺失手指", Type: "Negative", Category: "人物", Synonyms: []string{}, UsageCount: 280},
			{Value: "extra digit", Label: "多余手指", Type: "Negative", Category: "人物", Synonyms: []string{}, UsageCount: 260},
		},
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonData)
}

func (s *Seeder) seedCategories() error {
	fmt.Println("Seeding categories...")
	
	categories := []struct {
		name     string
		catType  string
		parentID uint
	}{
		// 一级分类
		{"人物", "ATOM", 0},
		{"场景", "ATOM", 0},
		{"风格", "ATOM", 0},
		{"质量", "ATOM", 0},
		{"光照", "ATOM", 0},
		{"道具", "ATOM", 0},
		{"预设库", "PRESET", 0},
	}

	for _, cat := range categories {
		_, err := s.categoryService.CreateCategory(cat.name, cat.catType, cat.parentID)
		if err != nil {
			fmt.Printf("Warning: failed to create category %s: %v\n", cat.name, err)
		} else {
			fmt.Printf("  Created category: %s\n", cat.name)
		}
	}

	// 创建子分类
	var personCategory models.Category
	if err := s.db.Where("name = ? AND type = ?", "人物", "ATOM").First(&personCategory).Error; err != nil {
		return err
	}

	subCategories := []struct {
		name string
	}{
		{"发型"},
		{"眼睛"},
		{"服装"},
		{"姿势"},
		{"表情"},
	}

	for _, sub := range subCategories {
		_, err := s.categoryService.CreateCategory(sub.name, "ATOM", personCategory.ID)
		if err != nil {
			fmt.Printf("Warning: failed to create subcategory %s: %v\n", sub.name, err)
		} else {
			fmt.Printf("  Created subcategory: %s\n", sub.name)
		}
	}

	return nil
}

func (s *Seeder) seedAtoms() error {
	fmt.Println("Seeding atoms...")

	// Get category IDs
	var categories []models.Category
	if err := s.db.Find(&categories).Error; err != nil {
		return err
	}

	categoryMap := make(map[string]uint)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	atoms := []struct {
		value      string
		label      string
		atomType   string
		category   string
		synonyms   []string
		usageCount int
	}{
		{"masterpiece", "杰作", "Positive", "质量", []string{"best quality"}, 500},
		{"best quality", "最佳质量", "Positive", "质量", []string{"high quality"}, 480},
		{"1girl", "1个女孩", "Positive", "人物", []string{"solo girl"}, 600},
		{"long hair", "长发", "Positive", "发型", []string{"flowing hair"}, 350},
		{"short hair", "短发", "Positive", "发型", []string{"bob cut"}, 280},
		{"blonde hair", "金发", "Positive", "发型", []string{"golden hair"}, 250},
		{"black hair", "黑发", "Positive", "发型", []string{"dark hair"}, 400},
		{"blue eyes", "蓝眼睛", "Positive", "眼睛", []string{"azure eyes"}, 300},
		{"red eyes", "红眼睛", "Positive", "眼睛", []string{"crimson eyes"}, 200},
		{"school uniform", "校服", "Positive", "服装", []string{"seifuku"}, 320},
		{"dress", "连衣裙", "Positive", "服装", []string{"gown"}, 280},
		{"standing", "站立", "Positive", "姿势", []string{"stand"}, 350},
		{"sitting", "坐着", "Positive", "姿势", []string{"sit"}, 280},
		{"smile", "微笑", "Positive", "表情", []string{"grin"}, 400},
		{"anime style", "动漫风格", "Positive", "风格", []string{"anime"}, 450},
		{"realistic", "写实风格", "Positive", "风格", []string{"photorealistic"}, 350},
		{"low quality", "低质量", "Negative", "质量", []string{"bad quality"}, 500},
		{"bad anatomy", "糟糕的人体结构", "Negative", "人物", []string{}, 400},
		{"bad hands", "糟糕的手", "Negative", "人物", []string{}, 380},
		{"text", "文字", "Negative", "质量", []string{}, 350},
		{"watermark", "水印", "Negative", "质量", []string{}, 320},
		{"blurry", "模糊", "Negative", "质量", []string{}, 300},
	}

	for _, atom := range atoms {
		categoryID := categoryMap[atom.category]
		if categoryID == 0 {
			fmt.Printf("Warning: category not found for %s\n", atom.category)
			continue
		}

		_, err := s.atomService.CreateAtom(atom.value, atom.label, atom.atomType, categoryID, atom.synonyms)
		if err != nil {
			fmt.Printf("Warning: failed to create atom %s: %v\n", atom.value, err)
		} else {
			fmt.Printf("  Created atom: %s (%s)\n", atom.value, atom.label)
		}
	}

	return nil
}

func (s *Seeder) seedPresets() error {
	fmt.Println("Seeding presets...")

	presets := []struct {
		title   string
		posText string
		negText string
		params  map[string]interface{}
	}{
		{
			title:   "动漫女孩基础预设",
			posText: "masterpiece, best quality, 1girl, solo, long hair, blue eyes, smile, school uniform, anime style",
			negText: "low quality, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, blurry",
			params: map[string]interface{}{
				"steps":   30,
				"cfg":     7.0,
				"sampler": "DPM++ 2M Karras",
				"model":   "animeModel_v20.safetensors",
				"width":   512,
				"height":  768,
			},
		},
		{
			title:   "写实风景预设",
			posText: "masterpiece, best quality, landscape, mountain, lake, sunset, golden hour, realistic, detailed, 8k uhd, photorealistic",
			negText: "low quality, blurry, cartoon, anime, painting, drawing, sketch, watermark, text, signature",
			params: map[string]interface{}{
				"steps":   35,
				"cfg":     7.5,
				"sampler": "DPM++ 2M SDE Karras",
				"model":   "realisticVisionV51_v51VAE.safetensors",
				"width":   1024,
				"height":  576,
			},
		},
		{
			title:   "赛博朋克风格",
			posText: "cyberpunk, neon lights, cityscape, night, rain, futuristic, sci-fi, detailed, cinematic lighting, 8k, highly detailed",
			negText: "low quality, blurry, daytime, sunny, historical, medieval, natural, organic, blurry",
			params: map[string]interface{}{
				"steps":   28,
				"cfg":     8.0,
				"sampler": "DPM++ SDE Karras",
				"model":   "deliberate_v3.safetensors",
				"width":   1024,
				"height":  576,
			},
		},
	}

	for _, preset := range presets {
		_, err := s.presetService.CreatePreset(
			preset.title,
			0, // categoryID - default to 0 for seed data
			preset.posText,
			preset.negText,
			[]uint{}, // atom IDs
			preset.params,
			[]string{}, // previews
		)
		if err != nil {
			fmt.Printf("Warning: failed to create preset %s: %v\n", preset.title, err)
		} else {
			fmt.Printf("  Created preset: %s\n", preset.title)
		}
	}

	return nil
}
