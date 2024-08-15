package tools

import (
	data "fs_scan/data"
	model "fs_scan/model"
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

func handleFile(file fs.DirEntry, db *data.FileDb, path string) {
	info, err := file.Info()
	if err != nil {
		return
	}
	if !contains(file.Name(), []string{`.jpg`, `.png`, `.jpeg`, `.pdf`}) {
		return
	}
	createdAt := info.ModTime()

	q := model.FileModel{
		Path:    filepath.Join(path, info.Name()),
		Format:  filepath.Ext(info.Name()),
		ModTime: createdAt,
		Tags:    []string{""},
	}
	db.AddQueryToGroup(q)
}

func scanFS(path string, fileHandler func(fs.DirEntry, string)) {
	var dirWg sync.WaitGroup
	worker(path, &dirWg, fileHandler)
	dirWg.Wait()
}

func ReloadDb(db *data.FileDb) {
	db.Exec("DELETE FROM files")
	var readerWg sync.WaitGroup
	readerWg.Add(1)
	go func() {
		defer readerWg.Done()
		db.StartInsertGroupingManager()
	}()
	scanFS("/mnt/c/Users/sergey/", func(de fs.DirEntry, s string) { handleFile(de, db, s) })
	db.StopGroupManager()
	readerWg.Wait()
}
