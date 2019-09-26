package mylogger

import (
	"fmt"
	"os"
	"time"
)

// ConsoleLogger 记录日志到终端
type ConsoleLogger struct {
	Level Level
}

// NewConsoleLogger ConsoleLogger的构造函数
func NewConsoleLogger(Level Level) *ConsoleLogger {
	return &ConsoleLogger{
		Level: Level,
	}
}

// WriteConsoleLogger 往终端控制台打印日志
func (c *ConsoleLogger) WriteConsoleLogger(level Level, format string, args ...interface{}) {
	if c.Level > level {
		return
	}

	msg := fmt.Sprintf(format, args...)
	now := time.Now().Format("2006-01-02 03:04:05")
	funcname, filename, line := GetCallerInfo(2)
	levelstr := ConvertLevelTOLevelstring(level)
	logmsg := fmt.Sprintf(
		"[%s][%s:%d][%s][%s]%s",
		now,
		filename,
		line,
		funcname,
		levelstr,
		msg,
	)
	fmt.Fprintln(os.Stdout, logmsg)
}

// Debug 打印debug日志
func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	c.WriteConsoleLogger(DebugLevel, format, args...)
}

// Info 打印debug日志
func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	c.WriteConsoleLogger(InfoLevel, format, args...)
}

// Warning 打印warning日志
func (c *ConsoleLogger) Warning(format string, args ...interface{}) {
	c.WriteConsoleLogger(WarningLevel, format, args...)
}

// Error 打印error日志
func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	c.WriteConsoleLogger(ErrorLevel, format, args...)
}

// Fatal 打印Fatal日志
func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	c.WriteConsoleLogger(FatalLevel, format, args...)
}

// Close 操作系统的标准输出不需要关闭，这里定义该方法完全是为了满足logger这个接口的定义。
func (c *ConsoleLogger) Close() {

}
