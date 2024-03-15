package internal

import (
	"time"
)

type DB struct {
	Table map[string]Entry
}

func InitializeDB() DB {
	return DB{
		Table: make(map[string]Entry),
	}
}

type Entry struct {
	Value string
	PX    int64
}

func (db DB) GetValue(key string) (string, bool) {
	entry, exists := db.Table[key]
	if !exists {
		return "", false
	}
	if entry.PX < time.Now().UnixMilli() {
		delete(db.Table, entry.Value)
		return "", false
	}
	return entry.Value, true
}

func (db DB) SetValue(key string, entry Entry) {
	db.Table[key] = entry
}
