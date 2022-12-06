package models

import "time"

// Model of search form on home page
type Search struct {
	Id        int
	UserId    int
	Query     string
	CreatedAt time.Time
}
