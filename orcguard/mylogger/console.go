package mylogger

import (
	"fmt"
	"os"
	"time"
)

type ConsoleLogger struct {
	level Level
}

func NewConsoleLogger(levelstr string) *ConsoleLogger {
	// 即: info --> InfoLevel
	loglevel := parseLogLevel(levelstr)
	cl := &ConsoleLogger{
		level: loglevel,
	}
	return cl
}

func (c *ConsoleLogger) log(level Level, format string, args ...interface{}) {
	if c.level > level {
		return
	}

	msg := fmt.Sprintf(format, args...)
	nowStr := time.Now().Format("2006-01-02 15:04:05.000")
	fileName, line, funcName := getCallerInfo(3)
	logLevelStr := getLevelStr(level) // 将InfoLevel 转成 INFO
	logMsg := fmt.Sprintf("[%s] [%s:%d] [%s] [%s] %s",
		nowStr, fileName, line, funcName, logLevelStr, msg)
	fmt.Fprintln(os.Stdout, logMsg)

}

func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	c.log(DebugLevel, format, args...)
}

func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	c.log(InfoLevel, format, args...)
}

func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	c.log(WarningLevel, format, args...)
}

func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	c.log(ErrorLevel, format, args...)
}

func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	c.log(FatalLevel, format, args...)
}

func (c *ConsoleLogger) Close() {

}
