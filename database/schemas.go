package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)


func createTables(db *sql.DB) error {
	
	// Create videos table
	videoTable := `
	CREATE TABLE IF NOT EXISTS videos (
    video_id VARCHAR(50) UNIQUE PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP NOT NULL,
    channel_title VARCHAR(255) NOT NULL,
    thumbnail_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(videoTable)
	if err != nil {
		return err
	}

	return nil
}

func MigrateUp() error {
	
	err := createTables(DB)
	if err != nil {
		return err
	}

	return nil
}