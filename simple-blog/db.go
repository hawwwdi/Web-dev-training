package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sql", "hadi:123%Ha%123@tcp(127.0.0.1)/blog")
	handleErr(os.Stderr, err)
	err = db.Ping()
	handleErr(os.Stderr, err)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}
