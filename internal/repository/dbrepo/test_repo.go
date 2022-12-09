package dbrepo

import (
	"errors"
	"time"

	"github.com/mcsymiv/web-hello-world/internal/models"
)

func (p *testDbRepo) GetUserSearchByUserIdAndFullTextQuery(userId int, s string) (models.Search, error) {
	var m models.Search = models.Search{
		Id:          1,
		UserId:      1,
		Query:       "test query",
		Description: "test description",
		Link:        "http://mcs.com?test=1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return m, nil
}

func (p *testDbRepo) GetUserSearchesByUserIdAndPartialTextQuery(userId int, s string) ([]models.Search, error) {
	var ms []models.Search

	if userId == 5 {
		return ms, errors.New("test error from db")
	}

	ms = []models.Search{
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

	return ms, nil
}

// InsertSearch inserts search entry to database
func (p *testDbRepo) InsertSearch(s models.Search) error {
	return nil
}
