package main

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "hadi:123%Ha%123@tcp(127.0.0.1)/blog")
	runtime.SetFinalizer(db, func(obj *sql.DB) {
		handleErr(os.Stderr, obj.Close())
		fmt.Println("database disconnected :*")
	})
	handleErr(os.Stderr, err)
	err = db.Ping()
	handleErr(os.Stderr, err)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	fmt.Println("data base connected :)")
}
