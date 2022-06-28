package logic

import (
	"fmt"
	"orcguard/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db       *sqlx.DB
	username = "orchestrator"
	password = "orch_monitorpd"
	dbname   = ""
)

func InitDB(ip string, port int) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True", username, password, ip, port, dbname)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect oldmaster DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

func get_readonly() (val string) {
	sqlStr := "show variables like 'read_only'"
	var m models.Mysql_vars
	err := db.Get(&m, sqlStr)
	if err != nil {
		fmt.Printf("get readonly failed, err:%v\n", err)
		return
	}
	fmt.Printf("name:%s val:%s\n", m.Name, m.Value)
	return m.Value
}
