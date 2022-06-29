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
	dbname   = ""
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

func GetReadOnly() (v string) {
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

func KillConnection() {
	sqlStr := " select id from information_schema.processlist where db is not null"
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	defer rows.Close()

	var IDs []int
	for rows.Next() {
		var num int
		err := rows.Scan(&num)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		IDs = append(IDs, num)
	}
	L.Info("id: %v", IDs)
	if len(IDs) >= 1 {
		KillId(IDs)
	}

}

func KillId(IDs []int) {
	for _, id := range IDs {
		sqlStr := "kill ?"
		db.QueryRow(sqlStr, id).Scan()
	}
	L.Info("ids killed")
}
