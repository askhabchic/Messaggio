package models

type Message struct {
	ID          int    `json:"id"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
	Status      bool   `json:"status"`
	ProcessedAt string `json:"processed_at,omitempty"`
}
