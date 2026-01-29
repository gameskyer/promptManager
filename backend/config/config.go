package config

import (
	"os"
	"path/filepath"
)

const (
	AppName    = "PromptMaster"
	AppVersion = "2.0.0"
)

var (
	AppDataDir string
	DBPath     string
	BackupDir  string
	ImageDir   string
)

func init() {
	// 使用当前工作目录下的数据库文件
	DBPath = "./promptmaster.db"
	
	// 应用数据目录
	homeDir, _ := os.UserHomeDir()
	AppDataDir = filepath.Join(homeDir, ".promptmaster")
	BackupDir = filepath.Join(AppDataDir, "backups")
	ImageDir = filepath.Join(AppDataDir, "images")
}

// EnsureAppDataDir ensures the application data directory exists
func EnsureAppDataDir() {
	// 确保数据库所在目录存在（当前目录总是存在，但以防万一）
	os.MkdirAll(".", 0755)
	os.MkdirAll(AppDataDir, 0755)
	os.MkdirAll(BackupDir, 0755)
	os.MkdirAll(ImageDir, 0755)
}

// GetImageDir returns the image storage directory
func GetImageDir() string {
	return ImageDir
}

// GetDBPath returns the database path
func GetDBPath() string {
	return DBPath
}
