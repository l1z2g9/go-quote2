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


func GetDB() *sql.DB {
    var db *sql.DB
    db, err := sql.Open("postgres", "host=ec2-23-21-189-181.compute-1.amazonaws.com port=5432 user=rvxfododvigsgj password=fb71ac2a80b2877ba7b4a3114fa4c6b4fa8bdaba6c34a15f0faf6edc07ddbec3 dbname=d9flf74p9len8l sslmode=disable")
    if err != nil {
		log.Fatal(err, "Fail to connect DB")
		Error.Fatal(err, "Fail to connect DB")
	}

Info.Println("Get connection")
	return db
}