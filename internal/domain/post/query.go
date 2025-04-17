package post

import "github.com/google/uuid"

type SearchQuery struct {
	Query   string    `json:"query;omitempty"`
	Tags    []string  `json:"tags;omitempty"`
	Title   string    `json:"title;omitempty"`
	Content string    `json:"content;omitempty"`
	UserID  uuid.UUID `json:"user_id;omitempty"`
}
