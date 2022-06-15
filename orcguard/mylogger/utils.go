package mylogger

import (
	"path"
	"runtime"
)

// getCallerInfo (getFileName, line, funcName)
func getCallerInfo(skip int) (fileName string, line int, funcName string) {
	pc, fileName, line, ok := runtime.Caller(skip)
	if !ok {
		return
	}
	fileName = path.Base(fileName)
	funcName = runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName)
	return
}
