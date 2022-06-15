package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

type FileLogger struct {
	level     Level
	fileName  string
	filepath  string
	file      *os.File
	errorFile *os.File
	maxsize   int64
}

func NewFileLogger(levelstr, filePath, fileName string) *FileLogger {
	// mylogger.NewFileLogger("info")  --> InfoLevel
	// 即: info --> InfoLevel
	loglevel := parseLogLevel(levelstr)
	fl := &FileLogger{
		level:    loglevel,
		fileName: fileName,
		filepath: filePath,
		maxsize:  100 * 1024 * 1024, // 100MB切日志
	}
	fl.initFile()
	return fl
}

func (f *FileLogger) initFile() {
	logName := path.Join(f.filepath, f.fileName)
	fileObj, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("open file error: %s %v", logName, err))
	}
	f.file = fileObj

	errLogName := fmt.Sprintf("%s.err", logName)
	errFileObj, err := os.OpenFile(errLogName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("open errfile error: %s %v", errLogName, err))
	}
	f.errorFile = errFileObj
}

func (f *FileLogger) log(level Level, format string, args ...interface{}) {
	// f.level = mylogger.NewFileLogger("info")
	// level = logger.Info == InfoLevel 常量名: InfoLevel 值: 1
	if f.level > level {
		return
	}

	msg := fmt.Sprintf(format, args...)
	nowStr := time.Now().Format("2006-01-02 15:04:05.000")
	fileName, line, funcName := getCallerInfo(3)
	logLevelStr := getLevelStr(level)
	logMsg := fmt.Sprintf("[%s] [%s:%d] [%s] [%s] %s", nowStr, fileName, line, funcName, logLevelStr, msg)
	// split log
	if f.checkSplit(f.file) {
		f.file = f.splitLogFile(f.file)
	}
	// write log
	fmt.Fprintln(f.file, logMsg)

	if level >= ErrorLevel {
		// split errlog
		if f.checkSplit(f.errorFile) {
			f.errorFile = f.splitLogFile(f.errorFile)
		}
		// write errlog
		fmt.Fprintln(f.errorFile, logMsg)
	}
}

func (f *FileLogger) checkSplit(file *os.File) bool {
	fileInfo, _ := file.Stat()
	filesize := fileInfo.Size()
	if filesize >= f.maxsize {
		return true
	} else {
		return false
	}
}

func (f *FileLogger) splitLogFile(file *os.File) *os.File {
	fileName := file.Name()
	nowStr := time.Now().Format("2006-01-02_150405")
	bakupName := fmt.Sprintf("%s_%v", fileName, nowStr)
	os.Rename(fileName, bakupName)
	file.Close()
	fileObj, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("split log file error!"))
	}
	return fileObj
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	f.log(DebugLevel, format, args...)
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	f.log(InfoLevel, format, args...)
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	f.log(WarningLevel, format, args...)
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	f.log(ErrorLevel, format, args...)
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	f.log(FatalLevel, format, args...)
}

func (f *FileLogger) Close() {
	f.file.Close()
	f.errorFile.Close()
}
