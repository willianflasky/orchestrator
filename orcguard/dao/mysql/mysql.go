package mysql

import (
	"fmt"
	. "orcguard/mylogger"
)

func OpertionDB_dao(oldmaster, newmaster string) (RWDomain, RODomain string, err error) {
	tx, err := db.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			fmt.Println("rollback")
			L.Error("db result: rollback")
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			fmt.Println("commit")
			L.Info("db result: commit")
		}
	}()
	// 1. get domain_rw
	sqlStr1 := `select domain from info where ip = ? and isdel=0;`
	err = tx.Get(&RWDomain, sqlStr1, oldmaster)
	if err != nil {
		fmt.Printf("get failed, ip: [%v] ,err: %v\n", oldmaster, err)
		L.Error("get failed, ip: [%v] ,err: %v\n", oldmaster, err)
		return
	}

	// 2. get domain_ro
	sqlStr2 := `select domain from info where ip = ? and rw =0 and isdel=0;`
	err = tx.Get(&RODomain, sqlStr2, newmaster)
	if err != nil {
		fmt.Printf("get failed, ip: [%v] ,err: %v\n", newmaster, err)
		L.Error("get failed, ip: [%v] ,err: %v\n", newmaster, err)
		return
	}

	// 3. domain_rw update new ip
	sqlStr3 := `update info set ip = ? where domain= ? and isdel=0;`
	ret, err := tx.Exec(sqlStr3, newmaster, RWDomain)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		L.Error("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		L.Error("get RowsAffected failed, err:%v\n", err)
		return
	}

	if n != 1 {
		L.Error("RowsAffected: %v", n)
	} else {
		L.Info("sql: %v RowsAffected: %v", sqlStr2, n)
	}

	// 4. domain_ro how many A
	sqlStr4 := "select count(ip) from info where domain= ? and isdel=0;"
	ro_num := 0
	err = tx.Get(&ro_num, sqlStr4, RODomain)
	if err != nil {
		L.Error("get %v num failed ,err: %v\n", RODomain, err)
		return
	}
	// 5. domain_ro a record > 1
	if ro_num > 1 {
		//  delete A record from domain_ro.
		sqlStr5 := "update info set isdel = 1 where domain= ? and ip = ? and isdel=0"
		ret, err = tx.Exec(sqlStr5, RODomain, newmaster)
		if err != nil {
			fmt.Printf("update failed, err:%v\n", err)
			L.Error("update failed, err:%v\n", err)
			return
		}
		n, err = ret.RowsAffected() // 操作影响的行数
		if err != nil {
			L.Error("get RowsAffected failed, err:%v\n", err)
			return
		}
		if n != 1 {
			L.Error("RowsAffected: %v", n)
		} else {
			L.Info("sql: %v RowsAffected: %v", sqlStr5, n)
		}
	} else {
		L.Info("only one A record.")
	}

	return RWDomain, RODomain, err
}
