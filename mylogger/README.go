// 在其他包中引用日志库，例如在show/main.go引用该日志库的代码如下：
package main

import (
	"mylogger"
)

func main() {
	var logger mylogger.Logger
	// consolelogger := mylogger.NewConsoleLogger(mylogger.InfoLevel)
	// logger = consolelogger
	filelogger := mylogger.NewFileLogger(mylogger.InfoLevel, "info.log", "d:/")
	logger = filelogger
	logger.Close()

	// 用一个统一的接口名字调用这些方法，所以无论是console打印还是日志记录，下面的代码都不需要更改。
	for i := 0; i < 100000; i++ {
		logger.Debug("%s日志测试", "debug")
		logger.Info("%s日志测试", "info")
		logger.Warning("%s日志测试", "warning")
		logger.Error("%s日志测试", "error")
		logger.Fatal("%s日志测试", "fatal")
	}
}
