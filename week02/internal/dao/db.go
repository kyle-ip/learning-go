package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db = newDb()

func newDb() *sql.DB {
	aDB, err := sql.Open("mysql",
		"root:mypassword@tcp(127.0.0.1:3306)/testdb?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	if err = aDB.Ping(); err != nil {
		panic(err)
	}
	aDB.SetMaxIdleConns(10)
	aDB.SetMaxOpenConns(30)
	aDB.SetConnMaxLifetime(time.Second * 5)
	return aDB
}
