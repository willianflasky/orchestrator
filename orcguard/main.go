package main

import (
	"flag"
	"fmt"
	"orcguard/dao/mysql"
	"orcguard/logic"
	. "orcguard/mylogger"
	"os"
	"path/filepath"
)

var (
	path, _   = filepath.Abs(os.Args[0])
	base_dir  = filepath.Dir(path)
	logspath  = filepath.Join(base_dir, "logs")
	oldmaster string
	newmaster string
)

func main() {
	// mkdir -p logs
	if _, err := os.Stat(logspath); os.IsNotExist(err) {
		os.MkdirAll(logspath, os.ModePerm)
	}

	// init logger
	InitLogger(logspath)
	defer L.Close()

	// init database.
	if err := mysql.InitDB(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		L.Error("init mysql failed")
		return
	}
	defer mysql.Close()

	// args
	flag.StringVar(&oldmaster, "old", "", "old master. eg: 10.1.1.1")
	flag.StringVar(&newmaster, "new", "", "new master. eg: 10.1.1.2")
	flag.Parse()

	// 1. init Info and
	info := logic.NewInfo(oldmaster, newmaster)
	info.Run()
}
