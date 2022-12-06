package repository

import "github.com/mcsymiv/web-hello-world/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertSearch(s models.Search) error
}
