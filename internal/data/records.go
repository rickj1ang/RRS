package data

import "time"

// impl MarshalJSON()([]byte, err) method for futhur customize the json
type Record struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Writer      string    `json:"writer,omitempty"`
	TotalPages  uint16    `json:"total_pages,omitempty"`
	CurrentPage uint16    `json:"current_page,omitempty"`
	Progress    float32   `json:"progress,omitempty"`
	Description string    `json:"description,omitempty"`
	Genres      []string  `json:"genres,omitempty"`
}
