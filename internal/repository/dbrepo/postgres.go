package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/mcsymiv/web-hello-world/internal/models"
)

func (p *postgresDBRepo) GetUserSearch(s string, userId int) (models.Search, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	i := `
		select id, query, created_at
		from searches
		where query = $1
		and user_id = $2
	`
	r, err := p.DB.QueryContext(ctx, i, s, userId)
	if err != nil {
		log.Println("Unable to get query for user", err)
		return models.Search{}, err
	}

	var id int
	var q string
	var t time.Time

	for r.Next() {
		err = r.Scan(&id, &q, &t)
		if err != nil {
			log.Println("Unable to get search")
			return models.Search{}, err
		}
	}

	return models.Search{
		Id:        id,
		UserId:    userId,
		Query:     q,
		CreatedAt: t,
	}, nil
}

func (p *postgresDBRepo) InsertSearch(s models.Search) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	i := `
		insert into 
		searches 
		(query, user_id, created_at)
		values 
		($1, $2, $3)
		`

	p.DB.ExecContext(ctx, i, s.Query, s.UserId, time.Now())

	return nil
}
