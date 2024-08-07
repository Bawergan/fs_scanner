package data

// import (
// 	"os"
// 	"sync"
// 	"testing"
// )

// const root = "./tmp_testing/"
// const dbFile = root + "test_db"

// func TestMain(m *testing.M) {
// 	exit := m.Run()
// 	os.RemoveAll(root)
// 	os.Exit(exit)
// }

// func TestOpenCloseDatabase(t *testing.T) {
// 	db, err := CreateFileDb(dbFile)
// 	defer os.Remove(dbFile)
// 	if err != nil {
// 		t.Errorf("OpenDatabase() returned an error: %v", err)
// 	}
// 	if db.db == nil {
// 		t.Error("db.db == nil after OpenDatabase()")
// 	}
// 	err = db.Close()
// 	if err != nil {
// 		t.Errorf("db.Close() returned an error: %v", err)
// 	}
// }
// func prepDb() (*FileDb, error) {
// 	db, err := CreateFileDb(dbFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sqlStmt := `
//         CREATE TABLE IF NOT EXISTS test_table (
//             id INTEGER PRIMARY KEY AUTOINCREMENT,
//             name TEXT NOT NULL
//         );
//     `
// 	err = db.CreateTable(sqlStmt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }
// func TestCreateTable(t *testing.T) {
// 	db, err := prepDb()
// 	if err != nil {
// 		t.Errorf("db prep failed: %s", err)
// 	}
// 	defer db.Close()
// 	defer os.Remove(dbFile)
// 	sqlStmt := `
//         CREATE TABLE test_table (
//             id INTEGER PRIMARY KEY AUTOINCREMENT,
//             name TEXT NOT NULL
//         );
//     `
// 	err = db.CreateTable(sqlStmt)
// 	if err == nil {
// 		t.Errorf("CreateTable() didn't return an error")
// 	}
// 	sqlStmt = `
//         CRETE TABLE test_table (
//             id IN NULL
//         );
//     `
// 	err = db.CreateTable(sqlStmt)
// 	if err == nil {
// 		t.Errorf("CreateTable() didn't return an error")
// 	}
// }

// func TestInsert(t *testing.T) {
// 	db, err := prepDb()
// 	if err != nil {
// 		t.Errorf("db prep failed: %s", err)
// 	}
// 	defer db.Close()
// 	defer os.Remove(dbFile)

// 	err = db.Insert(GenericQuery{query: "INSERT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
// 	if err != nil {
// 		t.Errorf("Insert() returned an error: %v", err)
// 	}
// 	err = db.Insert(GenericQuery{query: "INSRT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
// 	if err == nil {
// 		t.Errorf("Insert() didn't returne an error")
// 	}
// }

// func TestInsertGroup(t *testing.T) {
// 	db, err := prepDb()
// 	if err != nil {
// 		t.Errorf("db prep failed: %s", err)
// 	}
// 	defer db.Close()
// 	defer os.Remove(dbFile)
// 	var queryGroup []GenericQuery
// 	for range 50 {
// 		queryGroup = append(queryGroup, GenericQuery{query: "INSERT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
// 	}
// 	err = db.InsertGroup(queryGroup)
// 	if err != nil {
// 		t.Errorf("InsertGroup returned an error: %v", err)
// 	}
// 	queryGroup = append(queryGroup, GenericQuery{query: "asd"})
// 	err = db.InsertGroup(queryGroup)
// 	if err == nil {
// 		t.Errorf("InsertGroup returned an error: %v", err)
// 	}
// }

// func BenchmarkInsert(b *testing.B) {
// 	db, err := prepDb()
// 	if err != nil {
// 		b.Errorf("db prep failed: %s", err)
// 	}
// 	defer db.Close()
// 	defer os.Remove(dbFile)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		err = db.Insert(GenericQuery{query: "INSERT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
// 		if err != nil {
// 			b.Errorf("Insert() returned an error: %v", err)
// 		}
// 	}
// }
// func BenchmarkInsertGO(b *testing.B) {
// 	db, err := prepDb()
// 	if err != nil {
// 		b.Errorf("db prep failed: %s", err)
// 	}
// 	defer db.Close()
// 	defer os.Remove(dbFile)
// 	var wg sync.WaitGroup
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			db.Insert(GenericQuery{query: "INSERT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
// 		}()
// 	}
// 	wg.Wait()
// }

// func BenchmarkAutoGroupInsert(b *testing.B) {
// 	db, err := prepDb()
// 	if err != nil {
// 		b.Errorf("db prep failed: %s", err)
// 	}
// 	defer db.Close()
// 	defer os.Remove(dbFile)
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() { defer wg.Done(); db.StartInsertGroupingManager() }()
// 	b.ResetTimer()
// 	go func() {
// 		for i := 0; i < b.N; i++ {
// 			db.AddQueryToGroup(GenericQuery{query: "INSERT INTO test_table (name) VALUES (?)", args: []any{"test_name"}})
// 		}
// 	}()
// 	db.groupManagerLoop = false
// 	wg.Wait()
// }
