package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sql", "hadi:123%Ha%123@localhost/blog")
	handleErr(os.Stderr, err)
	err = db.Ping()
	handleErr(os.Stderr, err)
}



