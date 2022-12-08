package repository

import "github.com/mcsymiv/web-hello-world/internal/models"

type DatabaseRepo interface {
	GetUserSearchByUserIdAndFullTextQuery(userId int, s string) (models.Search, error)
	InsertSearch(s models.Search) error
}
