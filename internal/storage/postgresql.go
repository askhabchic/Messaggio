package storage

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("postgresql", storagePath)
	if err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}

}

// internal/database/postgres.go
package database

import (
"database/sql"
"project/internal/models"
_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}
	return db.Ping()
}

func SaveMessage(msg *models.Message) error {
	_, err := db.Exec("INSERT INTO messages (content, status) VALUES ($1, $2)", msg.Content, "pending")
	return err
}

func GetStats() (map[string]int, error) {
	rows, err := db.Query("SELECT status, COUNT(*) FROM messages GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats[status] = count
	}
	return stats, nil
}
