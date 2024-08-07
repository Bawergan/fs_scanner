package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	model "fs_scan/model"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const dbFilePath = `./store`
const dbName = `files`
const fileDbTable = `files (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        created_at DATETIME NOT NULL,
        tags TEXT NOT NULL
    )`

type FileDb struct {
	db               *sql.DB
	dbMu             sync.Mutex
	groupCh          chan model.FileModel
	groupManagerLoop bool
}

func CreateFileDb() (*FileDb, error) {
	err := os.MkdirAll(dbFilePath, 0755)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", filepath.Join(dbFilePath, dbName))
	if err != nil {
		return nil, err
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS ` + fileDbTable + `;`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return &FileDb{
		db:      db,
		groupCh: make(chan model.FileModel),
	}, nil
}

func (d *FileDb) Close() error {
	if d.db != nil {
		err := d.db.Close()
		return err
	}
	return fmt.Errorf("trying to close closed db")
}

func (d *FileDb) Insert(file model.FileModel) error {
	d.dbMu.Lock()
	defer d.dbMu.Unlock()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(model.FileModelQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	tags, err := json.Marshal(file.Tags)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(file.Path, file.ModTime, tags)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (d *FileDb) InsertGroup(fileGroup []model.FileModel) error {
	log.Println(len(fileGroup))

	d.dbMu.Lock()
	defer d.dbMu.Unlock()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, file := range fileGroup {

		stmt, err := tx.Prepare(model.FileModelQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		tags, err := json.Marshal(file.Tags)
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(file.Path, file.ModTime, tags)
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

func (d *FileDb) AddQueryToGroup(query model.FileModel) {
	d.groupCh <- query
}

func (d *FileDb) StopGroupManager() {
	d.groupManagerLoop = false
}

func (d *FileDb) StartInsertGroupingManager() {
	d.groupManagerLoop = true
	nextTurnOff := 1

	maxlen := 100000
	var group []model.FileModel
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
			// if len(group) == 0 && nextTurnOff != 10 {
			// 	time.Sleep(time.Millisecond * 100)
			// 	nextTurnOff += 1
			// } else {
			// 	d.groupManagerLoop = false
			// }
			if len(group) != 0 {
				err := d.InsertGroup(group)
				if err != nil {
					log.Fatal(err)
				}
				group = []model.FileModel{}
			}
		}
	}
}

func (d *FileDb) Query(query string) (*sql.Rows, error) {
	return d.db.Query(query)
}

func (d *FileDb) Exec(query string) (sql.Result, error) {
	return d.db.Exec(query)
}

func (d *FileDb) CountEnteries() (int, error) {
	rows, err := d.Query("SELECT count(*) FROM files")
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}

	err = rows.Err()
	if err != nil {
		return -1, err
	}
	return count, nil
}
