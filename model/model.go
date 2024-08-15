package model

import (
	"time"
)

type FileModel struct {
	Path    string
	Format  string
	ModTime time.Time
	Tags    []string
}

const FileModelQuery = "INSERT INTO files (name, string, created_at, tags) VALUES (?, ?, ?)"
