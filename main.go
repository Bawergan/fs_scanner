package main

import (
	"fmt"
	app "fs_scan/app"
	. "fs_scan/db"
	"log"
	"time"
)


func temp() {
	db, err := OpenDatabase("./files.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS files (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        created_at DATETIME NOT NULL,
        tags TEXT NOT NULL
    );
    `
	err = db.CreateTable(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	return

	rows, err := db.Query("SELECT name, created_at FROM files")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var createdAt time.Time
		err = rows.Scan(&name, &createdAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Name: %s, Created at: %s\n", name, createdAt)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := app.NewApp()
	app.LaunchCLI()
}
