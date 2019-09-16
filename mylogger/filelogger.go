package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// FileLogger 记录日志到文件的
type FileLogger struct {
	Level    Level
	FileName string
	FilePath string
	File     *os.File
	ErrFile  *os.File
	MaxSize  int64
}

// NewFileLogger 构造FileLogger实例
func NewFileLogger(Level Level, FileName, FilePath string) (f *FileLogger) {
	logfile := path.Join(FilePath, FileName)
	errLogfile := fmt.Sprintf("%s.error", logfile)
	File, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Errorf("打开文件%s 失败, 原因是：%s", logfile, err))
	}

	ErrFile, err := os.OpenFile(errLogfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Errorf("打开文件%s 失败, 原因是：%s", errLogfile, err))
	}
	return &FileLogger{
		Level:    Level,
		FileName: FileName,
		FilePath: FilePath,
		File:     File,
		ErrFile:  ErrFile,
		MaxSize:  1 * MB,
	}
}

// CheckSplit 是否切割
func (f *FileLogger) CheckSplit(file *os.File) bool {
	fileinfo, _ := file.Stat()
	filesize := fileinfo.Size()
	return filesize >= f.MaxSize
}

// SplitFile 切割日志
func (f *FileLogger) SplitFile(file *os.File) *os.File { // 传入的可能是正常日志文件，也可能是error日志文件句柄
	splittime := time.Now().Unix()
	filename := file.Name()                                  // 获取的filename是带路径的
	tarname := fmt.Sprintf("%s.log_%v", filename, splittime) // 拼接后的tarname也是带路径的。
	file.Close()
	os.Rename(filename, tarname)
	fileobj, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Errorf("Error： 打开文件%s失败", filename))
	}
	return fileobj
}

// WriterLog 记录日志
func (f *FileLogger) WriterLog(level Level, format string, args ...interface{}) {
	// FileLogger实例的level属性如果高于日志的某个级别，则不记录该级别的日志。
	if f.Level > level {
		/*
			WriterLog没有定义返回值，但是这里却有return，则实际上该函数不会反回任何数据，因此调用函数时
			也就不能使用变量接收，所以此处的return语句只起到提前结束函数执行的作用。
		*/
		return
	}

	msg := fmt.Sprintf(format, args...)
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	funcName, file, line := GetCallerInfo(2)
	levelstr := ConvertLevelTOLevelstring(level)

	// 生成要写入日志信息
	logmsg := fmt.Sprintf(
		"[%s][%s:%d][%s][%s]%s",
		nowStr,
		file,
		line,
		funcName,
		levelstr,
		msg,
	)

	if level >= ErrorLevel {
		if f.CheckSplit(f.ErrFile) {
			f.ErrFile = f.SplitFile(f.ErrFile)
		}
		fmt.Fprintln(f.ErrFile, logmsg)
		return
	}

	if f.CheckSplit(f.File) {
		f.File = f.SplitFile(f.File)
	}
	fmt.Fprintln(f.File, logmsg)
}

// Debug 输出Debug日志到文件
func (f *FileLogger) Debug(format string, args ...interface{}) {
	f.WriterLog(DebugLevel, format, args...)
}

// Info 记录info级别的日志。
func (f *FileLogger) Info(format string, args ...interface{}) {
	f.WriterLog(InfoLevel, format, args...)
}

// Warning 记录Warning级别的日志。
func (f *FileLogger) Warning(format string, args ...interface{}) {
	f.WriterLog(WarningLevel, format, args...)
}

// Error 记录Error级别的日志。
func (f *FileLogger) Error(format string, args ...interface{}) {
	f.WriterLog(ErrorLevel, format, args...)
}

// Fatal 记录重大错误日志
func (f *FileLogger) Fatal(format string, args ...interface{}) {
	f.WriterLog(FatalLevel, format, args...)
}

// Close 关闭结构体中保存的日志文件句柄
func (f *FileLogger) Close() {
	f.File.Close()
	f.ErrFile.Close()
}