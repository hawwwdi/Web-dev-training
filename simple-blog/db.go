package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"runtime"
	"time"
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
	creatTable()
}

func creatTable() {
	query := `CREATE TABLE IF NOT EXISTS users
	(
    	id       INT AUTO_INCREMENT PRIMARY KEY,
    	username VARCHAR(20) NOT NULL,
    	password VARCHAR(20) NOT NULL
	);`
	stmt, err := db.Prepare(query)
	handleErr(os.Stderr, err)
	defer stmt.Close()
	_, err1 := stmt.Exec()
	handleErr(os.Stderr, err1)
}
