package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// FileLogger 记录日志到文件的
type FileLogger struct {
	Level       Level
	FileName    string
	FilePath    string
	File        *os.File
	ErrFile     *os.File
	MaxSize     int64
	LogDataChan chan *LogData
}

// LogData 日志相关信息
type LogData struct {
	LogLevel string
	Message  string
	TimeStr  string
	FuncName string
	FileName string
	LineNo   int
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

	logger := &FileLogger{
		Level:    Level,
		FileName: FileName,
		FilePath: FilePath,
		File:     File,
		ErrFile:  ErrFile,
		MaxSize:  1 * MB,
		// 通道最多只能存5万条日志，超过的将丢弃
		LogDataChan: make(chan *LogData, 50000),
	}

	/*
		在构造函数中调用结构体实例的方法，由于此时构造函数还没有return，因为你不能像其他方法那样，
		使用f.WriteLogBackend()，实际上这个f就是logger这个结构体实例。一定要注意在logger结构体
		初始化完成后才能够启动下面的goroutine，因为该goroutine会读取logger实例中的LogDataChan
		通道，否则读取一个未初始化的引用类型（通道）会引发runtime error: invalid memory address
		or nil pointer dereference的panic
	*/
	go logger.WriteLogBackend()

	return logger
}

// CheckSplit 是否切割
func (f *FileLogger) CheckSplit(file *os.File) bool {
	fileinfo, err := file.Stat()
	if err != nil {
		panic(fmt.Errorf("获取日志文件属性错误，%s", err))
	}
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

// WriterLogToChan 发送日志到通道
func (f *FileLogger) WriterLogToChan(level Level, format string, args ...interface{}) {
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

	logdata := &LogData{
		LogLevel: levelstr,
		Message:  msg,
		TimeStr:  nowStr,
		FuncName: funcName,
		FileName: file,
		LineNo:   line,
	}

	// 通道满了之后会阻塞，这样会阻塞业务代码，这里采取的措施是丢弃最新的日志
	select {
	case f.LogDataChan <- logdata:
	// 什么都不做，表示丢弃最新的日志
	default:
	}
}

// WriteLogBackend 真正往文件里面写日志的函数
func (f *FileLogger) WriteLogBackend() {
	for logdata := range f.LogDataChan {
		// 生成要写入日志信息
		logmsg := fmt.Sprintf(
			"[%s][%s:%d][%s][%s]%s",
			logdata.TimeStr,
			logdata.FileName,
			logdata.LineNo,
			logdata.FuncName,
			logdata.LogLevel,
			logdata.Message,
		)

		if f.CheckSplit(f.File) {
			f.File = f.SplitFile(f.File)
		}
		fmt.Fprintln(f.File, logmsg)

		level := ConvertLevelstringTOLevel(logdata.LogLevel)
		if level >= ErrorLevel {
			if f.CheckSplit(f.ErrFile) {
				f.ErrFile = f.SplitFile(f.ErrFile)
			}
			fmt.Fprintln(f.ErrFile, logmsg)
		}
	}
}

// Debug 输出Debug日志到文件
func (f *FileLogger) Debug(format string, args ...interface{}) {
	f.WriterLogToChan(DebugLevel, format, args...)
}

// Info 记录info级别的日志。
func (f *FileLogger) Info(format string, args ...interface{}) {
	f.WriterLogToChan(InfoLevel, format, args...)
}

// Warning 记录Warning级别的日志。
func (f *FileLogger) Warning(format string, args ...interface{}) {
	f.WriterLogToChan(WarningLevel, format, args...)
}

// Error 记录Error级别的日志。
func (f *FileLogger) Error(format string, args ...interface{}) {
	f.WriterLogToChan(ErrorLevel, format, args...)
}

// Fatal 记录重大错误日志
func (f *FileLogger) Fatal(format string, args ...interface{}) {
	f.WriterLogToChan(FatalLevel, format, args...)
}

// Close 关闭结构体中保存的日志文件句柄
func (f *FileLogger) Close() {
	/*
		有条件的关闭日志文件句柄，因为还有一个并行的goroutine在往日志文件中写入日志，如果
		此时关闭文件，则该goroutine在执行checksplit函数时，由于要判断文件的大小，对一个已
		关闭了的文件句柄查询文件的属性，就会引发panic: runtime error: invalid memory
		address or nil pointer dereference异常。
	*/
	for {
		if len(f.LogDataChan) == 0 {
			f.File.Close()
			f.ErrFile.Close()
			break
		}
		time.Sleep(time.Millisecond)
	}
}
