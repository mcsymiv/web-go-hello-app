package models

import "time"

type User struct {
	Id        int
	UserName  string
	Email     string
	Password  string
	Searches  []Search
	CreatedAt time.Time
	UpdatedAt time.Time
}
