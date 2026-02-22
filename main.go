package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

// imageMiddleware 提供图片静态文件服务
func imageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查是否是图片请求
		if strings.HasPrefix(r.URL.Path, "/images/") {
			// 从 URL 中提取图片名称
			imageName := strings.TrimPrefix(r.URL.Path, "/images/")
			// 去除可能的反斜杠（Windows）
			imageName = strings.ReplaceAll(imageName, "\\", "/")
			imageName = filepath.Base(imageName)
			imagePath := filepath.Join(config.ImageDir, imageName)

			fmt.Printf("[DEBUG] Serving image: %s -> %s\n", r.URL.Path, imagePath)

			// 安全检查：确保路径在 ImageDir 内
			absPath, err := filepath.Abs(imagePath)
			if err != nil {
				http.Error(w, "Invalid path", http.StatusBadRequest)
				return
			}
			absImageDir, _ := filepath.Abs(config.ImageDir)
			if !strings.HasPrefix(absPath, absImageDir) {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}

			// 检查文件是否存在
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}

			// 提供文件
			http.ServeFile(w, r, imagePath)
			return
		}

		// 非图片请求，交给下一个 handler
		if next != nil {
			next.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
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
	batchService := services.NewBatchService(db)

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
	batchHandler := handlers.NewBatchHandler(batchService)

	// Create app
	app := NewApp()

	err = wails.Run(&options.App{
		Title:     "PromptMaster - AI绘画提示词管理",
		Width:     1200,
		Height:    750,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Middleware: imageMiddleware,
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
			batchHandler,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
