package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db       *sqlx.DB
	username = "orcguard"
	password = "orcguard"
	ip       = "10.0.9.168"
	port     = 61106
	dbname   = "orcguard"
)

func InitDB() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True", username, password, ip, port, dbname)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

func Close() {
	_ = db.Close()
}
