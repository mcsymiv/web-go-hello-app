package repository

import "github.com/mcsymiv/web-hello-world/internal/models"

type DatabaseRepo interface {
	GetUserSearch(s string, userId int) (models.Search, error)
	InsertSearch(s models.Search) error
}
