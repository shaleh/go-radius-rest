package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

func sqlConnect(sc *ServerConfig) (db *sql.DB, err error) {
	config := new(mysql.Config)
	config.User = sc.Database.User
	config.Passwd = sc.Database.Password
	config.DBName = sc.Database.Name
	config.Addr = sc.Database.Host
	config.Net = "tcp"

	db, err = sql.Open("mysql", config.FormatDSN())
	if err == nil {
		err = db.Ping() // actually open a connection to test the DSN and options
	}

	return
}

var (
	_db *sql.DB
)

func getDB() *sql.DB {
	return _db
}
