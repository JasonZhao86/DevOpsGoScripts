package mylogger

import (
	"path"
	"runtime"
)

// GetCallerInfo 获取调用者函数的信息
func GetCallerInfo(skip int) (funcName, file string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		funcPath := runtime.FuncForPC(pc).Name()
		funcName := path.Base(funcPath) // 获取到的函数带有包名，剥离出函数名
		filename := path.Base(file)   // 获取到的filename带路径，从file中剥离出文件名
		return funcName, filename, line
	}
	return // 返回string和int的默认零值
}