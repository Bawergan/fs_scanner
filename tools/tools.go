package tools

import (
	. "fs_scan/db"
	"io/fs"
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

func worker(root string, dirWg *sync.WaitGroup, fileHandler func(fs.DirEntry, string)) {
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
			worker(filepath.Join(root, ent.Name()), dirWg, fileHandler)
		} else {
			fileWg.Add(1)
			go func() {
				defer fileWg.Done()
				fileHandler(ent, root)
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
func scanFS(path string, fileHandler func(fs.DirEntry, string)) {
	var dirWg sync.WaitGroup
	worker(path, &dirWg, fileHandler)
	dirWg.Wait()
}
func ReloadDb(db *Database) {
	db.Exec("DELETE FROM files")
	var readerWg sync.WaitGroup
	readerWg.Add(1)
	go func() {
		defer readerWg.Done()
		db.StartInsertGroupingManager()
	}()
	scanFS("/", func(de fs.DirEntry, s string) { handleFile(de, db, s) })
	db.StopGroupManager()
	readerWg.Wait()
}
