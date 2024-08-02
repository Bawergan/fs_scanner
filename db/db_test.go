package db

import (
	"os"
	"sync"
	"testing"
)

func TestOpenCloseDatabase(t *testing.T) {
	db, err := OpenDatabase("./test")
	defer os.Remove("./test")
	if err != nil {
		t.Errorf("OpenDatabase() returned an error: %v", err)
	}
	if db.db == nil {
		t.Error("db.db == nil after OpenDatabase()")
	}
	err = db.Close()
	if err != nil {
		t.Errorf("db.Close() returned an error: %v", err)
	}
}

func TestCreateTable(t *testing.T) {
	db, err := OpenDatabase("./test")
	if err != nil {
		t.Errorf("OpenDatabase() returned an error: %v", err)
	}
	defer db.Close()
	defer os.Remove("./test")

	sqlStmt := `
        CREATE TABLE IF NOT EXISTS test_table (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        );
    `
	err = db.CreateTable(sqlStmt)
	if err != nil {
		t.Errorf("CreateTable() returned an error: %v", err)
	}

	sqlStmt = `
        CREATE TABLE test_table (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        );
    `
	err = db.CreateTable(sqlStmt)
	if err == nil {
		t.Errorf("CreateTable() didn't return an error")
	}
	sqlStmt = `
        CRETE TABLE test_table (
            id IN NULL
        );
    `
	err = db.CreateTable(sqlStmt)
	if err == nil {
		t.Errorf("CreateTable() didn't return an error")
	}
}

func TestInsert(t *testing.T) {
	db, err := OpenDatabase("./test")
	if err != nil {
		t.Errorf("OpenDatabase() returned an error: %v", err)
	}
	defer db.Close()
	defer os.Remove("./test")

	sqlStmt := `
        CREATE TABLE IF NOT EXISTS test_table (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        );
    `
	err = db.CreateTable(sqlStmt)
	if err != nil {
		t.Errorf("CreateTable() returned an error: %v", err)
	}

	err = db.Insert(GenericQuery{query: "INSERT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
	if err != nil {
		t.Errorf("Insert() returned an error: %v", err)
	}
	err = db.Insert(GenericQuery{query: "INSRT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
	if err == nil {
		t.Errorf("Insert() didn't returne an error")
	}
}

func TestInsertGroup(t *testing.T) {
	db, err := OpenDatabase("./test")
	if err != nil {
		t.Errorf("OpenDatabase() returned an error: %v", err)
	}
	defer db.Close()
	defer os.Remove("./test")

	sqlStmt := `
        CREATE TABLE IF NOT EXISTS test_table (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        );
    `
	err = db.CreateTable(sqlStmt)
	if err != nil {
		t.Errorf("CreateTable() returned an error: %v", err)
	}
	var queryGroup []GenericQuery
	for range 50 {
		queryGroup = append(queryGroup, GenericQuery{query: "INSERT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
	}
	err = db.InsertGroup(queryGroup)
	if err != nil {
		t.Errorf("InsertGroup returned an error: %v", err)
	}
	queryGroup = append(queryGroup, GenericQuery{query: "asd"})
	err = db.InsertGroup(queryGroup)
	if err == nil {
		t.Errorf("InsertGroup returned an error: %v", err)
	}
}

func BenchmarkInsert(b *testing.B) {
	db, err := OpenDatabase("./test")
	if err != nil {
		b.Errorf("OpenDatabase() returned an error: %v", err)
	}
	defer db.Close()
	defer os.Remove("./test")

	sqlStmt := `
        CREATE TABLE IF NOT EXISTS test_table (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        );
    `
	err = db.CreateTable(sqlStmt)
	if err != nil {
		b.Errorf("CreateTable() returned an error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = db.Insert(GenericQuery{query: "INSERT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
		if err != nil {
			b.Errorf("Insert() returned an error: %v", err)
		}
	}
}
func BenchmarkInsertGO(b *testing.B) {
	db, err := OpenDatabase("./test")
	if err != nil {
		b.Errorf("OpenDatabase() returned an error: %v", err)
	}
	defer db.Close()
	defer os.Remove("./test")

	sqlStmt := `
        CREATE TABLE IF NOT EXISTS test_table (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        );
    `
	err = db.CreateTable(sqlStmt)
	if err != nil {
		b.Errorf("CreateTable() returned an error: %v", err)
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db.Insert(GenericQuery{query: "INSERT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
		}()
	}
	wg.Wait()
}

func BenchmarkAutoGroupInsert(b *testing.B) {
	db, err := OpenDatabase("./test")
	if err != nil {
		b.Errorf("OpenDatabase() returned an error: %v", err)
	}
	defer db.Close()
	defer os.Remove("./test")

	sqlStmt := `
		        CREATE TABLE IF NOT EXISTS test_table (
		            id INTEGER PRIMARY KEY AUTOINCREMENT,
		            name TEXT NOT NULL
		        );
		    `
	err = db.CreateTable(sqlStmt)
	if err != nil {
		b.Errorf("CreateTable() returned an error: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); db.StartInsertGroupingManager() }()
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			db.AddQueryToGroup(GenericQuery{query: "INSERT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
		}
	}()
	db.groupManagerLoop = false
	wg.Wait()
}
