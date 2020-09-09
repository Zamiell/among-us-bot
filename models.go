package main

import (
	"database/sql"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DatabaseFilename = "database.sqlite3"
)

var db *sql.DB

// Models contains a list of interfaces representing database tables
type Models struct {
	PlayerList
	Players
}

// modelsInit opens a database connection based on the credentials in the ".env" file
func modelsInit() (*Models, error) {
	databasePath := path.Join(projectPath, DatabaseFilename)
	if v, err := sql.Open("sqlite3", databasePath); err != nil {
		return nil, err
	} else {
		db = v
	}

	// Create the model
	return &Models{}, nil
}

// Close exposes the ability to close the underlying database connection
func (*Models) Close() {
	db.Close()
}
