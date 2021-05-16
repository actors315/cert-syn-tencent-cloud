package db

import (
	"database/sql"
	"fmt"
)

type Conn struct {
	Dsn string
	Db  *sql.DB
}

var QcloudToolDb Conn

func (conn Conn) Update(sqlStr string, args ...interface{}) (affected int64, err error) {

	stmt, err := conn.Db.Prepare(sqlStr)
	if nil != err {
		fmt.Println("failed to prepare query", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if nil != err {
		fmt.Println("failed to exec query", err)
		return 0, err
	}

	affected, err = result.RowsAffected()
	return
}

func (conn Conn) Insert(sqlStr string, args ...interface{}) (lastInsertId int64, err error) {
	stmt, err := conn.Db.Prepare(sqlStr)
	if nil != err {
		fmt.Println("failed to prepare query", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if nil != err {
		fmt.Println("failed to exec query", err)
		return 0, err
	}

	lastInsertId, err = result.LastInsertId()
	return
}

func (conn Conn) Query(sqlStr string, args ...interface{}) (rows *sql.Rows, err error) {
	var db = QcloudToolDb.Db
	stmt, err := db.Prepare(sqlStr)
	if nil != err {
		fmt.Println("failed to prepare query:", err)
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(args...)
	if nil != err {
		fmt.Println("failed to query data:", err)
		return
	}

	return
}
