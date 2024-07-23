package storage

import (
	"Messaggio/internal/config"
	"Messaggio/internal/models"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func NewStorage(db *sql.DB, log *slog.Logger) *Storage {
	return &Storage{
		db:  nil,
		log: nil,
	}
}

func (s *Storage) CreateTable() error {
	const fn = "storage.postgresql.CreateTable()"

	stmt, err := s.db.Prepare(`CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at timestamptz DEFAULT statement_timestamp(),
    status BOOLEAN DEFAULT FALSE,
    processed_at timestamptz);`)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) SaveMessage(msg models.Message) error {
	const fn = "storage.postgresql.SaveMessage()"

	err := s.db.QueryRow("INSERT INTO messages (content) VALUES ($1)", msg.Content)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetStats() (map[string]int, error) {
	const fn = "storage.postgresql.GetStatus()"

	stats := make(map[string]int)
	var count int
	if err := s.db.QueryRow(
		"SELECT COUNT(*) FROM messages WHERE status = TRUE",
	).Scan(&count); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	stats["processed"] = count

	if err := s.db.QueryRow(
		"SELECT COUNT(*) FROM messages WHERE status = FALSE",
	).Scan(&count); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	stats["unprocessed"] = count

	return stats, nil
}

func Connection(cfg *config.Config) (*sql.DB, error) {
	const fn = "storage.postgresql.New()"

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.DbName, cfg.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return db, nil
}
