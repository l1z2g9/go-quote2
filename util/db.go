package util

import (
	"database/sql"
	// _ "github.com/mattn/go-sqlite3"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
)

var dbPath string

func GetExecutePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Fatal(err)
	}
	return dir
}

func SetDB_old(path ...string) {
	if path != nil {
		dbPath = path[0]
	} else {
		dir := GetExecutePath()
		dbPath = dir + "/../stock.db"
	}
}

// GetDB Get the DB connection
func GetDB_old() *sql.DB {
	// log.Println("DB_PATH ", dbPath)
	if len(dbPath) == 0 {
		SetDB_old()
	}

	var db *sql.DB
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err, "dbPath = ", dbPath)
		Error.Fatal(err, "dbPath = ", dbPath)
	}

	return db
}

func GetDB() *sql.DB {
	var db *sql.DB
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err, "Fail to connect DB")
	}

	fmt.Println("Successfully connected!")
	return db
}
