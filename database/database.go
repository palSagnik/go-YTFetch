package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/palSagnik/go-YTFetch.git/config"
)

var (
	DB  *sql.DB
	err error
)

// connecting to the database
func ConnectDB() error {

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", config.DB_HOST, config.DB_USER, config.DB_PASS, config.DB_NAME, config.DB_PORT, config.SSL_MODE)
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	fmt.Println("connected to database")
	return nil
}
