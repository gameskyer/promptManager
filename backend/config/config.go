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
	
	// 应用数据目录（使用项目相对路径）
	AppDataDir = "./data"
	BackupDir = filepath.Join(AppDataDir, "backups")
	ImageDir = "./images"
}

// EnsureAppDataDir ensures the application data directory exists
func EnsureAppDataDir() {
	// 确保项目相对路径的目录存在
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
