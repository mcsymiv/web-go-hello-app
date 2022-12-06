package dbrepo

import "github.com/mcsymiv/web-hello-world/internal/models"

func (p *postgresDBRepo) AllUsers() bool {
	return true
}

func (p *postgresDBRepo) InsertSearch(s models.Search) error {
	i := `
		insert into 
		searches 
		(query, user_id, created_at)
		values 
		($1, $2, $3)
		`

	p.DB.Exec(i, s.Query, s.UserId, s.CreatedAt)

	return nil
}
