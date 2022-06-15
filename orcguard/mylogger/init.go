package mylogger

var L Logger

func InitLogger(logspath string) {
	L = NewFileLogger("debug", logspath, "orcguard.log")
}
