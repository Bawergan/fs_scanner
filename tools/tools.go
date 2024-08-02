package tools

import (
	. "fs_scan/db"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func contains(s string, args []string) bool {
	for _, arg := range args {
		if strings.Contains(s, arg) {
			return true
		}
	}
	return false
}

func worker(db *Database, root string, dirWg *sync.WaitGroup) {
	if contains(root, []string{`/mnt/c/ProgramData`, `/mnt/c/Windows`}) {
		return
	}
	dirWg.Add(1)
	defer dirWg.Done()
	dirEnt, err := os.ReadDir(root)
	if err != nil {
		return
	}
	var fileWg sync.WaitGroup
	for _, ent := range dirEnt {
		if ent.Type() == fs.ModeSymlink {
			continue
		}
		if ent.IsDir() {
			go worker(db, filepath.Join(root, ent.Name()), dirWg)
		} else {
			fileWg.Add(1)
			go func() {
				defer fileWg.Done()
				handleFile(ent, db, root)
			}()
		}
	}
	fileWg.Wait()
}

func handleFile(file fs.DirEntry, db *Database, path string) {
	info, err := file.Info()
	if err != nil {
		return
	}
	createdAt := info.ModTime()

	q := FileInsertionQuery{
		Path:    filepath.Join(path, info.Name()),
		ModTime: createdAt,
		Tags:    []string{""},
	}
	db.AddQueryToGroup(q.ConvertToGeneric())
}

func ReloadDb(db *Database){

	var readerWg sync.WaitGroup
	go func() {
		readerWg.Add(1)
		defer readerWg.Done()
		db.StartInsertGroupingManager()
	}()
	var dirWg sync.WaitGroup
	worker(db, "/home", &dirWg)
	log.Println("waiting for workers")
	dirWg.Wait()
	log.Println("waiting for reader")
	db.StopGroupManager()
	readerWg.Wait()
}
