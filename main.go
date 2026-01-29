package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	
	"promptmaster/backend/config"
	"promptmaster/backend/handlers"
	"promptmaster/backend/models"
	"promptmaster/backend/services"
	"promptmaster/backend/utils"
)

//go:embed frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// OnStartup is called when the app starts
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
}

func main() {
	// Ensure app data directory exists
	config.EnsureAppDataDir()
	
	// Initialize database
	db, err := models.InitDB()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	
	// Run seeder to import default data
	seeder := utils.NewSeeder(db)
	if err := seeder.SeedAll(); err != nil {
		fmt.Printf("Warning: Failed to seed database: %v\n", err)
		// Don't exit, continue with empty database
	}
	
	// Initialize services
	atomService := services.NewAtomService(db)
	presetService := services.NewPresetService(db)
	categoryService := services.NewCategoryService(db)
	versionService := services.NewVersionService(db)
	searchService := services.NewSearchService(db)
	aiService := services.NewAIService(db)
	imageService := services.NewImageService(db)
	backupService := services.NewBackupService(db)
	
	// Initialize handlers
	atomHandler := handlers.NewAtomHandler(atomService)
	presetHandler := handlers.NewPresetHandler(presetService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	versionHandler := handlers.NewVersionHandler(versionService)
	searchHandler := handlers.NewSearchHandler(searchService)
	aiHandler := handlers.NewAIHandler(aiService)
	seederHandler := handlers.NewSeederHandler(seeder)
	imageHandler := handlers.NewImageHandler(imageService)
	backupHandler := handlers.NewBackupHandler(backupService)
	
	// Create app
	app := NewApp()
	
	err = wails.Run(&options.App{
		Title:     "PromptMaster - AI绘画提示词管理",
		Width:     1200,
		Height:    750,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.OnStartup,
		Bind: []interface{}{
			atomHandler,
			presetHandler,
			categoryHandler,
			versionHandler,
			searchHandler,
			aiHandler,
			seederHandler,
			imageHandler,
			backupHandler,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
