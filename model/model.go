package model

import (
	"time"
)

type FileModel struct {
	Path    string
	ModTime time.Time
	Tags    []string
}

const FileModelQuery = "INSERT INTO files (name, created_at, tags) VALUES (?, ?, ?)"
