package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const loggerKey = "logger"

// ANSI 颜色代码
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

type Logger struct {
	instance  *log.Logger
	requestId string
}

// NewLogger 创建一个新的 Logger 实例
func NewLogger(requestId string) *Logger {
	l := log.New(os.Stderr, "", 0)
	return &Logger{
		instance:  l,
		requestId: requestId,
	}
}

// GetLogger 从 context 中获取 Logger，如果不存在则创建一个默认的
func GetLogger(ctx *gin.Context) *Logger {
	if logger, exists := ctx.Get(loggerKey); exists {
		return logger.(*Logger)
	}
	return NewLogger("default-" + strings.ReplaceAll(uuid.New().String(), "-", ""))
}

// SetLogger 将 Logger 存入 context
func SetLogger(ctx *gin.Context, logger *Logger) {
	ctx.Set(loggerKey, logger)
}

// formatLog 格式化日志行，包括时间戳、requestId、日志级别和消息
func (l *Logger) formatLog(level, color, format string, v ...any) string {
	timestamp := time.Now().Format(time.DateTime)
	message := fmt.Sprintf(format, v...)
	return fmt.Sprintf("%s[%s] [%s] [%s] %s%s", color, level, l.requestId, timestamp, message, colorReset)
}

// Info 输出 Info 级别日志（绿色）
func (l *Logger) Info(format string, v ...any) {
	l.instance.Println(l.formatLog("INFO", colorGreen, format, v...))
}

// Error 输出 Error 级别日志（红色）
func (l *Logger) Error(format string, v ...any) {
	l.instance.Println(l.formatLog("ERROR", colorRed, format, v...))
}

// Warn 输出 Warn 级别日志（黄色）
func (l *Logger) Warn(format string, v ...any) {
	l.instance.Println(l.formatLog("WARN", colorYellow, format, v...))
}

// Debug 输出 Debug 级别日志（蓝色）
func (l *Logger) Debug(format string, v ...any) {
	l.instance.Println(l.formatLog("DEBUG", colorBlue, format, v...))
}
