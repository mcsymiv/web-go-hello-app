package repository

import "github.com/mcsymiv/web-hello-world/internal/models"

type DatabaseRepo interface {
	GetUserSearchByUserIdAndFullTextQuery(userId int, s string) (models.Search, error)
	GetUserSearchesByUserIdAndPartialTextQuery(userId int, s string) ([]models.Search, error)
	GetUserById(userId int) (models.User, error)
	GetUsersCount() (int, error)
	GetUsers() ([]models.User, error)
	GetSearchesByUserId(id int) ([]models.Search, error)
	GetUserSearchById(userId, searchId int) (models.Search, error)

	AuthenticateUser(email, password string) (int, string, error)
	InsertSearch(s models.Search) error

	AddUser(u models.User) error
	UpdateUser(u models.User) error
	UpdateUserSearch(s models.Search, userId int) error

	DeleteUserSearch(searchId, userId int) error
}
