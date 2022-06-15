# mylogger

1. main.go列出使用方式。
2. 当mylogger.NewFileLogger("info", "./", "a.log")初始化为info时，即小于info的日志，即：debug就不再打印了。
3. file.go -> NewFileLogger  -> maxsize 默认日志大于100MB就能切割。
