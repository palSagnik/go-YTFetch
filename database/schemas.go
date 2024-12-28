package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)


func createTables(db *sql.DB) error {
	return nil
}

func MigrateUp() error {
	
	err := createTables(DB)
	if err != nil {
		return err
	}

	return nil
}