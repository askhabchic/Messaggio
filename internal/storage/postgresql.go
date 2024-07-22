package storage

import (
	"Messaggio/internal/models"
	"database/sql"
	"fmt"
	//"project/internal/models"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) New(storagePath string) (*Storage, error) {
	const fn = "storage.postgresql.New()"
	db, err := sql.Open("postgresql", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at timestamptz DEFAULT statement_timestamp(),
    status BOOLEAN DEFAULT FALSE,
    processed_at timestamptz);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveMessage(msg *models.Message) error {
	const fn = "storage.postgresql.SaveMessage()"

	err := s.db.QueryRow("INSERT INTO messages (content) VALUES ($1)", msg.Content)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetStats() {
	const fn = "storage.postgresql.GetStatus()"

	s.db.Query("SELECT status, COUNT(*) FROM messages GROUP BY status")

}

//func GetStats() (map[string]int, error) {
//	rows, err := db.Query("SELECT status, COUNT(*) FROM messages GROUP BY status")
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	stats := make(map[string]int)
//	for rows.Next() {
//		var status string
//		var count int
//		if err := rows.Scan(&status, &count); err != nil {
//			return nil, err
//		}
//		stats[status] = count
//	}
//	return stats, nil
//}
