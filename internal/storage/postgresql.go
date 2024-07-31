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
		db:  db,
		log: log,
	}
}

const (
	tableCreateQuery = `CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at timestamptz DEFAULT statement_timestamp(),
    status BOOLEAN DEFAULT FALSE,
    processed_at timestamptz);`
	insertMessagesQuery    = `INSERT INTO messages (content) VALUES ($1)`
	selectStatusTrueQuery  = `SELECT COUNT(*) FROM messages WHERE status = TRUE`
	selectStatusFalseQuery = `SELECT COUNT(*) FROM messages WHERE status = FALSE`
)

func (s *Storage) CreateTable() error {
	const fn = "storage.postgresql.CreateTable()"

	_, err := s.db.Exec(tableCreateQuery)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) SaveMessage(msg models.Message) error {
	const fn = "storage.postgresql.SaveMessage()"

	_, err := s.db.Exec(insertMessagesQuery, msg.Content)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetStats() (map[string]int, error) {
	const fn = "storage.postgresql.GetStatus()"

	stats := make(map[string]int)
	var count int
	rows := s.db.QueryRow(selectStatusTrueQuery).Scan(&count)
	if rows.Error() != "" {
		return nil, fmt.Errorf("%s: %w", fn, rows.Error())
	}
	stats["processed"] = count

	rows = s.db.QueryRow(selectStatusFalseQuery).Scan(&count)
	if rows.Error() != "" {
		return nil, fmt.Errorf("%s: %w", fn, rows.Error())
	}
	stats["unprocessed"] = count

	return stats, nil
}

func Connection(cfg *config.Config) (*sql.DB, error) {
	const fn = "storage.postgresql.New()"

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DbHost, cfg.Username, cfg.Password, cfg.DbName, cfg.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return db, nil
}
