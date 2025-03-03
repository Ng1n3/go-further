package models

import (
	"errors"
	"fmt"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID         int
	Title      string
	Content    string
	Created_at time.Time
	Expires    time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created_at     time.Time
}

func (s *Snippet) String() string {
	return fmt.Sprintf("ID: %d\nTitle: %s\nContent: %s\nCreated: %s\nExpires: %s\n",
		s.ID, s.Title, s.Content, s.Created_at.Format(time.RFC3339), s.Expires.Format(time.RFC3339))
}
