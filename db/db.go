package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type FileInsertionQuery struct {
	Path    string
	ModTime time.Time
	Tags    []string
}

func (q FileInsertionQuery) ConvertToGeneric() GenericQuery {
	tags, err := json.Marshal(q.Tags)
	if err != nil {
		log.Fatal(err)
	}
	return GenericQuery{query: "INSERT INTO files (name, created_at, tags) VALUES (?, ?, ?)", args: []any{q.Path, q.ModTime, tags}}
}

type GenericQuery struct {
	query string
	args  []any
}

type Database struct {
	db               *sql.DB
	dbMu             sync.Mutex
	groupCh          chan GenericQuery
	groupManagerLoop bool
}

func OpenDatabase(filename string) (*Database, error) {
	log.Println("db: OpenDatabase called")

	db, err := sql.Open("sqlite3", filename+".db")
	if err != nil {
		return nil, err
	}

	return &Database{
		db:      db,
		groupCh: make(chan GenericQuery),
	}, nil
}

func (d *Database) Close() error {
	log.Println("db: db.Close called")

	if d.db != nil {
		err := d.db.Close()
		return err
	}

	return fmt.Errorf("trying to close closed db")
}

func (d *Database) CreateTable(sqlStmt string) error {
	log.Println("db: db.CreateTable called")

	_, err := d.db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) Insert(query GenericQuery) error {
	d.dbMu.Lock()
	defer d.dbMu.Unlock()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query.query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(query.args...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) InsertGroup(queryGroup []GenericQuery) error {
	log.Println(len(queryGroup))
	d.dbMu.Lock()
	defer d.dbMu.Unlock()
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, query := range queryGroup {

		stmt, err := tx.Prepare(query.query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(query.args...)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) AddQueryToGroup(query GenericQuery) {
	d.groupCh <- query
}
func (d *Database) StopGroupManager() {
	d.groupManagerLoop = false
}
func (d *Database) StartInsertGroupingManager() {
	d.groupManagerLoop = true
	nextTurnOff := 1

	maxlen := 10000
	var group []GenericQuery
	for d.groupManagerLoop || len(group) > 0 {
		//log.Println(d.groupManagerLoop, len(group), nextTurnOff)

		for len(group) >= maxlen/nextTurnOff {
			err := d.InsertGroup(group[:maxlen/nextTurnOff])
			if err != nil {
				log.Fatal(err)
			}
			group = group[maxlen/nextTurnOff:]
		}
		select {
		case q := <-d.groupCh:
			group = append(group, q)
		case <-time.After(time.Millisecond * 30):
			if len(group) == 0 && nextTurnOff != 10 {
				time.Sleep(time.Millisecond * 500)
				nextTurnOff += 1
			} else {
				d.groupManagerLoop = false
			}
			err := d.InsertGroup(group)
			if err != nil {
				log.Fatal(err)
			}
			group = []GenericQuery{}
		}
	}
}
func (d *Database) Query(query string) (*sql.Rows, error) {
	return d.db.Query(query)
}
