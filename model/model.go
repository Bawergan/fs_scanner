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

const FileModelQuery = "INSERT INTO files (name, format, created_at, tags) VALUES (?, ?, ?, ?)"
