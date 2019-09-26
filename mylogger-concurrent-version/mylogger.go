package mylogger

import (
	"strings"
)

// Level 定义日志级别类型
type Level uint16

// Logger 定义一个统一的接口
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Close()
}

const (
	// DebugLevel 字符串类型的日志级别比数字更见名知意
	DebugLevel Level = iota
	// InfoLevel info级别
	InfoLevel
	// WarningLevel 警告级别
	WarningLevel
	// ErrorLevel 错误级别
	ErrorLevel
	// FatalLevel 重大错误
	FatalLevel
)

const (
	_ = iota
	// KB = 1024
	KB = 1 << (iota * 10)
	// MB 左移20位
	MB
	// GB 左移30位
	GB
	// TB 左移40位
	TB
	// PB 左移50位
	PB
)

// ConvertLevelTOLevelstring 将日志级别转换成对应的字符串
func ConvertLevelTOLevelstring(level Level) string {
	switch level {
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case WarningLevel:
		return "Warning"
	case ErrorLevel:
		return "Error"
	case FatalLevel:
		return "Fatal"
	default:
		return "Info"
	}
}

// ConvertLevelstringTOLevel 将字符串格式的日志级别转换成对应的日志级别类型。
func ConvertLevelstringTOLevel(level string) Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warning":
		return WarningLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}