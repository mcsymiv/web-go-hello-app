package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/mcsymiv/web-hello-world/internal/models"
)

func (p *postgresDBRepo) GetUserSearchByUserIdAndFullTextQuery(userId int, s string) (models.Search, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	i := `
		select id, query, description, link, created_at, updated_at
		from searches
		where user_id = $1
		and query = $2
	`

	r, err := p.DB.QueryContext(ctx, i, userId, s)
	if err != nil {
		log.Println("Unable to get query for user", err)
		return models.Search{}, err
	}

	var m models.Search

	for r.Next() {
		err = r.Scan(&m.Id, &m.Query, &m.Description, &m.Link, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			log.Println("Unable to get search")
			return models.Search{}, err
		}
	}

	return m, nil
}

// InsertSearch inserts search entry to database
func (p *postgresDBRepo) InsertSearch(s models.Search) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	i := `
		insert into 
		searches 
		(query, user_id, link, description, created_at, updated_at)
		values 
		($1, $2, $3, $4, $5, $6)
	`

	p.DB.ExecContext(ctx, i, s.Query, s.UserId, s.Link, s.Description, time.Now(), time.Now())

	return nil
}
