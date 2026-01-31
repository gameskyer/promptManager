package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger 提供日志记录功能，支持文件分割
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	filename string
	maxSize  int64 // 最大文件大小（字节）
	size     int64 // 当前文件大小
}

// NewLogger 创建一个新的日志记录器
func NewLogger(filename string, maxSizeMB int) (*Logger, error) {
	if maxSizeMB <= 0 {
		maxSizeMB = 10 // 默认10MB
	}
	
	logger := &Logger{
		filename: filename,
		maxSize:  int64(maxSizeMB) * 1024 * 1024, // 转换为字节
	}
	
	// 确保日志目录存在
	logDir := filepath.Dir(filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	
	// 打开或创建日志文件
	if err := logger.openFile(); err != nil {
		return nil, err
	}
	
	return logger, nil
}

// openFile 打开或创建日志文件
func (l *Logger) openFile() error {
	file, err := os.OpenFile(l.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	
	// 获取当前文件大小
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to stat log file: %w", err)
	}
	
	l.file = file
	l.size = info.Size()
	
	return nil
}

// rotate 轮转日志文件
func (l *Logger) rotate() error {
	// 关闭当前文件
	if l.file != nil {
		l.file.Close()
	}
	
	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(l.filename)
	base := l.filename[:len(l.filename)-len(ext)]
	backupName := fmt.Sprintf("%s_%s%s", base, timestamp, ext)
	
	// 重命名当前文件
	if err := os.Rename(l.filename, backupName); err != nil {
		// 如果重命名失败，尝试删除旧文件
		os.Remove(l.filename)
	}
	
	// 打开新文件
	return l.openFile()
}

// Write 写入日志
func (l *Logger) Write(level string, message string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	// 检查是否需要轮转
	if l.size >= l.maxSize {
		if err := l.rotate(); err != nil {
			return err
		}
	}
	
	// 格式化日志行
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)
	
	// 写入文件
	n, err := l.file.WriteString(logLine)
	if err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}
	
	// 立即刷新到磁盘
	l.file.Sync()
	
	l.size += int64(n)
	return nil
}

// Info 记录信息日志
func (l *Logger) Info(message string) error {
	return l.Write("INFO", message)
}

// Debug 记录调试日志
func (l *Logger) Debug(message string) error {
	return l.Write("DEBUG", message)
}

// Error 记录错误日志
func (l *Logger) Error(message string) error {
	return l.Write("ERROR", message)
}

// LogAIRequest 记录 AI 请求日志
func (l *Logger) LogAIRequest(provider, endpoint, model string, requestBody interface{}) error {
	requestJSON, _ := json.MarshalIndent(requestBody, "", "  ")
	message := fmt.Sprintf("\n=== AI Request ===\nProvider: %s\nEndpoint: %s\nModel: %s\nRequest Body:\n%s\n", 
		provider, endpoint, model, string(requestJSON))
	return l.Info(message)
}

// LogAIResponse 记录 AI 响应日志
func (l *Logger) LogAIResponse(statusCode int, responseBody string, duration time.Duration) error {
	message := fmt.Sprintf("\n=== AI Response ===\nStatus: %d\nDuration: %s\nResponse Body:\n%s\n", 
		statusCode, duration, responseBody)
	return l.Info(message)
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// GetLogDir 返回日志目录（当前目录下的 logs 文件夹）
func GetLogDir() string {
	return "./logs"
}

// GetAILogPath 返回 AI 日志文件路径
func GetAILogPath() string {
	return filepath.Join(GetLogDir(), "ai_service.log")
}
