package dbrepo

import (
	"errors"
	"strings"
	"time"

	"github.com/mcsymiv/web-hello-world/internal/models"
)

var testSearchModels []models.Search = []models.Search{
	models.Search{
		Id:          1,
		UserId:      1,
		Query:       "test query",
		Description: "test description",
		Link:        "http://mcs.com?test=1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	models.Search{
		Id:          2,
		UserId:      2,
		Query:       "test query 2",
		Description: "test description 2",
		Link:        "http://mcs.com?test=2",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

func (p *testDbRepo) GetUserSearchByUserIdAndFullTextQuery(userId int, s string) (models.Search, error) {

	for _, m := range testSearchModels {
		if m.Query == s {
			return m, nil
		}
	}

	return testSearchModels[0], nil
}

func (p *testDbRepo) GetUserSearchesByUserIdAndPartialTextQuery(userId int, s string) ([]models.Search, error) {
	var ms []models.Search

	if userId == 5 {
		return ms, errors.New("test error from db")
	}

	for _, m := range testSearchModels {
		if strings.Contains(m.Query, s) {
			ms = append(ms, m)
		}
	}

	return ms, nil
}

// InsertSearch inserts search entry to database
func (p *testDbRepo) InsertSearch(s models.Search) error {
	if s.Query == "invalid" {
		return errors.New("test error from db on insert invalid 'search_query'")
	}

	testSearchModels = append(testSearchModels, s)
	return nil
}

// GetUserById retrieves full user data from db by id
func (p *testDbRepo) GetUserById(userId int) (models.User, error) {
	var u models.User
	return u, nil
}

// UpdateUser updates user info by id
func (p *testDbRepo) UpdateUser(u models.User) error {
	return nil
}

// AuthenticateUser compares user emain and pasword hash
func (p *testDbRepo) AuthenticateUser(email, password string) (int, string, error) {

	var id int
	var hash string

	return id, hash, nil
}

// GetUsersCount returns users count
func (p *testDbRepo) GetUsersCount() (int, error) {
	return 5, nil
}
