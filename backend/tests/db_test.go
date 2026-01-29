package tests

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"promptmaster/backend/models"
	"promptmaster/backend/services"
	"promptmaster/backend/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Create temp database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate
	err = db.AutoMigrate(
		&models.Category{},
		&models.Atom{},
		&models.Preset{},
		&models.PresetVersion{},
		&models.Preview{},
		&models.UsageStat{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestCategoryService(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewCategoryService(db)

	t.Run("CreateCategory", func(t *testing.T) {
		category, err := service.CreateCategory("测试分类", "ATOM", 0)
		if err != nil {
			t.Errorf("Failed to create category: %v", err)
		}
		if category.Name != "测试分类" {
			t.Errorf("Expected name '测试分类', got '%s'", category.Name)
		}
		if category.Type != "ATOM" {
			t.Errorf("Expected type 'ATOM', got '%s'", category.Type)
		}
	})

	t.Run("GetCategoryByID", func(t *testing.T) {
		// Create a category first
		created, _ := service.CreateCategory("查询测试", "ATOM", 0)

		// Query it
		category, err := service.GetCategoryByID(created.ID)
		if err != nil {
			t.Errorf("Failed to get category: %v", err)
		}
		if category.ID != created.ID {
			t.Errorf("Expected ID %d, got %d", created.ID, category.ID)
		}
	})

	t.Run("GetCategoriesByParent", func(t *testing.T) {
		// Create parent
		parent, _ := service.CreateCategory("父分类", "ATOM", 0)
		
		// Create children
		service.CreateCategory("子分类1", "ATOM", parent.ID)
		service.CreateCategory("子分类2", "ATOM", parent.ID)

		// Query children
		children, err := service.GetCategoriesByParent(parent.ID, "")
		if err != nil {
			t.Errorf("Failed to get children: %v", err)
		}
		if len(children) != 2 {
			t.Errorf("Expected 2 children, got %d", len(children))
		}
	})

	t.Run("UpdateCategory", func(t *testing.T) {
		category, _ := service.CreateCategory("更新前", "ATOM", 0)

		updated, err := service.UpdateCategory(category.ID, map[string]interface{}{
			"name": "更新后",
		})
		if err != nil {
			t.Errorf("Failed to update category: %v", err)
		}
		if updated.Name != "更新后" {
			t.Errorf("Expected name '更新后', got '%s'", updated.Name)
		}
	})

	t.Run("DeleteCategory", func(t *testing.T) {
		category, _ := service.CreateCategory("待删除", "ATOM", 0)

		err := service.DeleteCategory(category.ID)
		if err != nil {
			t.Errorf("Failed to delete category: %v", err)
		}

		// Try to find it
		_, err = service.GetCategoryByID(category.ID)
		if err == nil {
			t.Error("Expected error when getting deleted category, got nil")
		}
	})
}

func TestAtomService(t *testing.T) {
	db := setupTestDB(t)
	categoryService := services.NewCategoryService(db)
	atomService := services.NewAtomService(db)

	// Create test category
	category, _ := categoryService.CreateCategory("原子词测试分类", "ATOM", 0)

	t.Run("CreateAtom", func(t *testing.T) {
		atom, err := atomService.CreateAtom("test_value", "测试值", "Positive", category.ID, []string{"synonym1"})
		if err != nil {
			t.Errorf("Failed to create atom: %v", err)
		}
		if atom.Value != "test_value" {
			t.Errorf("Expected value 'test_value', got '%s'", atom.Value)
		}
		if atom.Label != "测试值" {
			t.Errorf("Expected label '测试值', got '%s'", atom.Label)
		}
	})

	t.Run("GetAtomByID", func(t *testing.T) {
		created, _ := atomService.CreateAtom("query_test", "查询测试", "Positive", category.ID, nil)

		atom, err := atomService.GetAtomByID(created.ID)
		if err != nil {
			t.Errorf("Failed to get atom: %v", err)
		}
		if atom.ID != created.ID {
			t.Errorf("Expected ID %d, got %d", created.ID, atom.ID)
		}
	})

	t.Run("GetAtomsByCategory", func(t *testing.T) {
		// Create multiple atoms
		for i := 0; i < 5; i++ {
			atomService.CreateAtom(
				"category_test_"+string(rune('0'+i)),
				"分类测试"+string(rune('0'+i)),
				"Positive",
				category.ID,
				nil,
			)
		}

		atoms, total, err := atomService.GetAtomsByCategory(category.ID, 1, 10)
		if err != nil {
			t.Errorf("Failed to get atoms by category: %v", err)
		}
		if total < 5 {
			t.Errorf("Expected at least 5 atoms, got %d", total)
		}
		if len(atoms) < 5 {
			t.Errorf("Expected at least 5 atoms in result, got %d", len(atoms))
		}
	})

	t.Run("FindAtomsBySynonym", func(t *testing.T) {
		// Create atom with synonyms
		atomService.CreateAtom("synonym_test", "同义词测试", "Positive", category.ID, []string{"test", "example"})

		// Search by synonym
		atoms, err := atomService.FindAtomsBySynonym("test")
		if err != nil {
			t.Errorf("Failed to find atoms by synonym: %v", err)
		}
		if len(atoms) == 0 {
			t.Error("Expected to find atoms, got none")
		}
	})

	t.Run("UpdateAtom", func(t *testing.T) {
		atom, _ := atomService.CreateAtom("update_test", "更新测试", "Positive", category.ID, nil)

		updated, err := atomService.UpdateAtom(atom.ID, map[string]interface{}{
			"label": "已更新",
		})
		if err != nil {
			t.Errorf("Failed to update atom: %v", err)
		}
		if updated.Label != "已更新" {
			t.Errorf("Expected label '已更新', got '%s'", updated.Label)
		}
	})

	t.Run("RecordUsage", func(t *testing.T) {
		atom, _ := atomService.CreateAtom("usage_test", "使用测试", "Positive", category.ID, nil)
		initialCount := atom.UsageCount

		err := atomService.RecordUsage(atom.ID)
		if err != nil {
			t.Errorf("Failed to record usage: %v", err)
		}

		// Verify usage count increased
		updated, _ := atomService.GetAtomByID(atom.ID)
		if updated.UsageCount != initialCount+1 {
			t.Errorf("Expected usage count %d, got %d", initialCount+1, updated.UsageCount)
		}
	})

	t.Run("DeleteAtom", func(t *testing.T) {
		atom, _ := atomService.CreateAtom("delete_test", "删除测试", "Positive", category.ID, nil)

		err := atomService.DeleteAtom(atom.ID)
		if err != nil {
			t.Errorf("Failed to delete atom: %v", err)
		}

		// Try to find it
		_, err = atomService.GetAtomByID(atom.ID)
		if err == nil {
			t.Error("Expected error when getting deleted atom, got nil")
		}
	})

	t.Run("BatchImportAtoms", func(t *testing.T) {
		jsonData := `[
			{"value": "import1", "label": "导入1", "type": "Positive", "category_id": 1, "synonyms": ["imp1"]},
			{"value": "import2", "label": "导入2", "type": "Positive", "category_id": 1, "synonyms": ["imp2"]}
		]`

		count, err := atomService.BatchImportAtoms(jsonData)
		if err != nil {
			t.Errorf("Failed to batch import atoms: %v", err)
		}
		if count != 2 {
			t.Errorf("Expected to import 2 atoms, got %d", count)
		}
	})
}

func TestPresetService(t *testing.T) {
	db := setupTestDB(t)
	presetService := services.NewPresetService(db)

	t.Run("CreatePreset", func(t *testing.T) {
		params := map[string]interface{}{
			"steps":   30,
			"cfg":     7.0,
			"sampler": "Euler a",
		}

		preset, err := presetService.CreatePreset(
			"测试预设",
			"masterpiece, best quality, 1girl",
			"low quality, bad anatomy",
			[]uint{},
			params,
		)

		if err != nil {
			t.Errorf("Failed to create preset: %v", err)
		}
		if preset.Title != "测试预设" {
			t.Errorf("Expected title '测试预设', got '%s'", preset.Title)
		}
		if preset.CurrentVersion != 1 {
			t.Errorf("Expected current version 1, got %d", preset.CurrentVersion)
		}
	})

	t.Run("GetPresetByID", func(t *testing.T) {
		created, _ := presetService.CreatePreset(
			"查询测试",
			"positive prompt",
			"negative prompt",
			[]uint{},
			map[string]interface{}{},
		)

		preset, err := presetService.GetPresetByID(created.ID)
		if err != nil {
			t.Errorf("Failed to get preset: %v", err)
		}
		if preset.ID != created.ID {
			t.Errorf("Expected ID %d, got %d", created.ID, preset.ID)
		}
	})

	t.Run("GetPresets", func(t *testing.T) {
		// Create multiple presets
		for i := 0; i < 3; i++ {
			presetService.CreatePreset(
				"预设"+string(rune('0'+i)),
				"positive",
				"negative",
				[]uint{},
				map[string]interface{}{},
			)
		}

		presets, total, err := presetService.GetPresets(1, 10, false)
		if err != nil {
			t.Errorf("Failed to get presets: %v", err)
		}
		if total < 3 {
			t.Errorf("Expected at least 3 presets, got %d", total)
		}
		if len(presets) < 3 {
			t.Errorf("Expected at least 3 presets in result, got %d", len(presets))
		}
	})

	t.Run("SoftDeletePreset", func(t *testing.T) {
		preset, _ := presetService.CreatePreset(
			"待删除",
			"positive",
			"negative",
			[]uint{},
			map[string]interface{}{},
		)

		err := presetService.SoftDeletePreset(preset.ID)
		if err != nil {
			t.Errorf("Failed to soft delete preset: %v", err)
		}

		// Get without deleted
		presets, _, _ := presetService.GetPresets(1, 10, false)
		for _, p := range presets {
			if p.ID == preset.ID {
				t.Error("Found deleted preset in non-deleted query")
			}
		}

		// Get with deleted
		allPresets, _, _ := presetService.GetPresets(1, 10, true)
		found := false
		for _, p := range allPresets {
			if p.ID == preset.ID && p.IsDeleted {
				found = true
				break
			}
		}
		if !found {
			t.Error("Deleted preset not found in include-deleted query")
		}
	})

	t.Run("RestorePreset", func(t *testing.T) {
		preset, _ := presetService.CreatePreset(
			"待恢复",
			"positive",
			"negative",
			[]uint{},
			map[string]interface{}{},
		)

		// Delete then restore
		presetService.SoftDeletePreset(preset.ID)
		err := presetService.RestorePreset(preset.ID)
		if err != nil {
			t.Errorf("Failed to restore preset: %v", err)
		}

		// Verify restored
		restored, _ := presetService.GetPresetByID(preset.ID)
		if restored.IsDeleted {
			t.Error("Preset should not be deleted after restore")
		}
	})

	t.Run("BuildPromptText", func(t *testing.T) {
		// Note: This requires atoms to exist in database
		// For this test, we'll use empty atom IDs
		text, ids, err := presetService.BuildPromptText([]uint{})
		if err != nil {
			t.Errorf("Failed to build prompt text: %v", err)
		}
		if text != "" {
			t.Errorf("Expected empty text, got '%s'", text)
		}
		if len(ids) != 0 {
			t.Errorf("Expected 0 IDs, got %d", len(ids))
		}
	})
}

func TestVersionService(t *testing.T) {
	db := setupTestDB(t)
	presetService := services.NewPresetService(db)
	versionService := services.NewVersionService(db)

	// Create a preset first
	preset, _ := presetService.CreatePreset(
		"版本测试",
		"positive",
		"negative",
		[]uint{},
		map[string]interface{}{"steps": 30},
	)

	t.Run("CreateVersion", func(t *testing.T) {
		input := services.VersionInput{
			PresetID: preset.ID,
			PosText:  "updated positive",
			NegText:  "updated negative",
			AtomIDs:  []uint{},
			Params:   map[string]interface{}{"steps": 35},
		}

		version, err := versionService.CreateVersion(input)
		if err != nil {
			t.Errorf("Failed to create version: %v", err)
		}
		if version.VersionNum != 2 {
			t.Errorf("Expected version num 2, got %d", version.VersionNum)
		}
		if version.PresetID != preset.ID {
			t.Errorf("Expected preset ID %d, got %d", preset.ID, version.PresetID)
		}
	})

	t.Run("GetVersion", func(t *testing.T) {
		version, err := versionService.GetVersion(preset.ID, 1)
		if err != nil {
			t.Errorf("Failed to get version: %v", err)
		}
		if version.VersionNum != 1 {
			t.Errorf("Expected version 1, got %d", version.VersionNum)
		}
	})

	t.Run("GetVersionHistory", func(t *testing.T) {
		// Create a few more versions
		for i := 0; i < 3; i++ {
			versionService.CreateVersion(services.VersionInput{
				PresetID: preset.ID,
				PosText:  "version " + string(rune('0'+i)),
				NegText:  "negative",
				AtomIDs:  []uint{},
				Params:   map[string]interface{}{},
			})
		}

		versions, err := versionService.GetVersionHistory(preset.ID, 10)
		if err != nil {
			t.Errorf("Failed to get version history: %v", err)
		}
		if len(versions) < 4 { // Original + 3 new
			t.Errorf("Expected at least 4 versions, got %d", len(versions))
		}
	})

	t.Run("StarVersion", func(t *testing.T) {
		// Get a version
		versions, _ := versionService.GetVersionHistory(preset.ID, 1)
		if len(versions) == 0 {
			t.Skip("No versions to star")
		}

		version := versions[0]
		err := versionService.StarVersion(version.ID, true)
		if err != nil {
			t.Errorf("Failed to star version: %v", err)
		}

		// Verify starred
		starred, _ := versionService.GetStarredVersions(preset.ID)
		found := false
		for _, v := range starred {
			if v.ID == version.ID && v.IsStarred {
				found = true
				break
			}
		}
		if !found {
			t.Error("Version not found in starred versions")
		}
	})

	t.Run("CompareVersions", func(t *testing.T) {
		diff, err := versionService.CompareVersions(preset.ID, 1, 2)
		if err != nil {
			t.Errorf("Failed to compare versions: %v", err)
		}
		if diff == nil {
			t.Error("Expected diff result, got nil")
		}
	})
}

func TestSeeder(t *testing.T) {
	db := setupTestDB(t)
	seeder := utils.NewSeeder(db)

	t.Run("SeedAll", func(t *testing.T) {
		err := seeder.SeedAll()
		if err != nil {
			t.Errorf("Failed to seed database: %v", err)
		}

		// Verify categories
		var categoryCount int64
		db.Model(&models.Category{}).Count(&categoryCount)
		if categoryCount == 0 {
			t.Error("Expected categories after seeding, got none")
		}

		// Verify atoms
		var atomCount int64
		db.Model(&models.Atom{}).Count(&atomCount)
		if atomCount == 0 {
			t.Error("Expected atoms after seeding, got none")
		}

		// Verify presets
		var presetCount int64
		db.Model(&models.Preset{}).Count(&presetCount)
		if presetCount == 0 {
			t.Error("Expected presets after seeding, got none")
		}
	})

	t.Run("SeedFromJSON", func(t *testing.T) {
		jsonData := `{
			"categories": [
				{"name": "JSON分类1", "type": "ATOM", "parent_id": 0},
				{"name": "JSON分类2", "type": "ATOM", "parent_id": 0}
			],
			"atoms": [
				{"value": "json_atom1", "label": "JSON原子1", "type": "Positive", "category": "JSON分类1", "synonyms": []},
				{"value": "json_atom2", "label": "JSON原子2", "type": "Positive", "category": "JSON分类2", "synonyms": []}
			]
		}`

		err := seeder.SeedFromJSON(jsonData)
		if err != nil {
			t.Errorf("Failed to seed from JSON: %v", err)
		}

		// Verify new category
		var category models.Category
		err = db.Where("name = ?", "JSON分类1").First(&category).Error
		if err != nil {
			t.Error("Expected to find JSON category, got error")
		}
	})

	t.Run("GetDefaultSeedData", func(t *testing.T) {
		data := seeder.GetDefaultSeedData()
		if data == "" {
			t.Error("Expected seed data, got empty string")
		}

		// Verify it's valid JSON
		var seedData utils.SeedData
		err := json.Unmarshal([]byte(data), &seedData)
		if err != nil {
			t.Errorf("Expected valid JSON, got error: %v", err)
		}
	})
}

func TestIntegration(t *testing.T) {
	db := setupTestDB(t)
	seeder := utils.NewSeeder(db)

	// Run full seed
	if err := seeder.SeedAll(); err != nil {
		t.Fatalf("Failed to seed for integration test: %v", err)
	}

	t.Run("EndToEndWorkflow", func(t *testing.T) {
		// 1. Verify categories exist
		var categories []models.Category
		if err := db.Where("type = ?", "ATOM").Find(&categories).Error; err != nil {
			t.Errorf("Failed to query categories: %v", err)
		}
		if len(categories) == 0 {
			t.Error("Expected categories in database")
		}

		// 2. Verify atoms exist
		var atoms []models.Atom
		if err := db.Find(&atoms).Error; err != nil {
			t.Errorf("Failed to query atoms: %v", err)
		}
		if len(atoms) == 0 {
			t.Error("Expected atoms in database")
		}

		// 3. Verify atoms have correct categories
		for _, atom := range atoms {
			if atom.CategoryID == 0 {
				t.Errorf("Atom %s has no category", atom.Value)
			}
		}

		// 4. Verify presets exist with versions
		var presets []models.Preset
		if err := db.Preload("Versions").Find(&presets).Error; err != nil {
			t.Errorf("Failed to query presets: %v", err)
		}
		if len(presets) == 0 {
			t.Error("Expected presets in database")
		}

		// 5. Verify each preset has at least one version
		for _, preset := range presets {
			if len(preset.Versions) == 0 {
				t.Errorf("Preset %s has no versions", preset.Title)
			}
		}
	})
}

// Benchmark tests
func BenchmarkCreateAtom(b *testing.B) {
	db := setupTestDB(&testing.T{})
	categoryService := services.NewCategoryService(db)
	atomService := services.NewAtomService(db)

	category, _ := categoryService.CreateCategory("Benchmark", "ATOM", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		atomService.CreateAtom(
			"bench_value_"+string(rune(i)),
			"基准测试",
			"Positive",
			category.ID,
			nil,
		)
	}
}

func BenchmarkGetAtomsByCategory(b *testing.B) {
	db := setupTestDB(&testing.T{})
	categoryService := services.NewCategoryService(db)
	atomService := services.NewAtomService(db)

	category, _ := categoryService.CreateCategory("Benchmark", "ATOM", 0)

	// Create 100 atoms
	for i := 0; i < 100; i++ {
		atomService.CreateAtom("value_"+string(rune(i)), "标签", "Positive", category.ID, nil)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		atomService.GetAtomsByCategory(category.ID, 1, 50)
	}
}
