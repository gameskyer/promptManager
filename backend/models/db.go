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
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	
	// Initialize default categories
	initDefaultCategories(db)
	
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
