package models

import "time"

// Model of search form on home page
type Search struct {
	Id          int
	UserId      int
	Query       string
	Link        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SplitSearch model for destructuring search query
type SplitSearch struct {
	Id        int
	Base      string
	Query     string
	Link      string
	Desc      string
	UpdatedAt time.Time
}
