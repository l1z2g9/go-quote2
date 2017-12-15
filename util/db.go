package util

import (
	"database/sql"
	// _ "github.com/mattn/go-sqlite3"
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

func SetDB(path ...string) {
	if path != nil {
		dbPath = path[0]
	} else {
		dir := GetExecutePath()
		dbPath = dir + "/../stock.db"
	}
}

// GetDB Get the DB connection
func GetDB() *sql.DB {
	// log.Println("DB_PATH ", dbPath)
	if len(dbPath) == 0 {
		SetDB()
	}

	var db *sql.DB
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err, "dbPath = ", dbPath)
		Error.Fatal(err, "dbPath = ", dbPath)
	}

	return db
}
