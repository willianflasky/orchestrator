package logic

import (
	"database/sql"
	"fmt"
	"orcguard/models"
	. "orcguard/mylogger"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db       *sql.DB
	username = "orchestrator"
	password = "orch_monitorpd"
	dbname   = "mysql"
)

func InitDB(ip string, port int) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", username, password, ip, port, dbname)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("connect oldmaster DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

func get_readonly() (val string) {
	sqlStr := `show variables like 'read_only'`
	var m models.Mysql_vars
	err := db.QueryRow(sqlStr).Scan(&m.Variable_name, &m.Value)
	if err != nil {
		fmt.Printf("get readonly failed, err:%v\n", err)
		return
	}
	L.Info("name:%s val:%s\n", m.Variable_name, m.Value)
	return m.Value
}
