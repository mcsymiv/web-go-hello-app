package repository

import "github.com/mcsymiv/web-hello-world/internal/models"

type DatabaseRepo interface {
	GetUserSearchByUserIdAndFullTextQuery(userId int, s string) (models.Search, error)
	GetUserSearchesByUserIdAndPartialTextQuery(userId int, s string) ([]models.Search, error)
	InsertSearch(s models.Search) error
	GetUserById(userId int) (models.User, error)
	UpdateUser(u models.User) error
	AuthenticateUser(email, password string) (int, string, error)
}
