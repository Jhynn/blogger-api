package database

import (
	"blogger/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Driver
)

// Connect returns a connection with the database - mysql (remember to close the connection).
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.ConnectionStringDB)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()

		return nil, err
	}

	return db, nil
}
