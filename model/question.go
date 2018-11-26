package model

import "time"

// Question interface
type Question struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	AuthorID int       `json:"author_id"`
	HasBest  bool      `json:"has_best"`
	Created  time.Time `json:"created"`
}
