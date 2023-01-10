package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mcsymiv/web-hello-world/internal/models"
	"golang.org/x/crypto/bcrypt"
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
			return m, err
		}
	}

	return m, nil
}

func (p *postgresDBRepo) GetUserSearchesByUserIdAndPartialTextQuery(userId int, s string) ([]models.Search, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	i := `
		select id, query, description, link, created_at, updated_at
		from searches
		where user_id = $1
		and query like '%' || $2 || '%'
	`

	var ms = make([]models.Search, 0)

	r, err := p.DB.QueryContext(ctx, i, userId, s)
	if err != nil {
		p.App.ErrorLog.Println("unable to get queries with text like for user", err)
		return ms, err
	}

	for r.Next() {
		var m models.Search
		err = r.Scan(&m.Id, &m.Query, &m.Description, &m.Link, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			log.Println("Unable to get search")
			return ms, err
		}

		ms = append(ms, m)
	}

	return ms, nil
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

// GetUserById retrieves full user data from db by id
func (p *postgresDBRepo) GetUserById(userId int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	i := `
		select id, user_name, email, password, searches, created_at, update_at
		from users
		where id = $1
		`
	row := p.DB.QueryRowContext(ctx, i, userId)

	var u models.User

	err := row.Scan(
		&u.Id,
		&u.UserName,
		&u.Email,
		&u.Password,
		&u.Searches,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		log.Println("Unable to get user")
		return u, err
	}

	return u, nil
}

// UpdateUser updates user info by id
func (p *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	i := `
		update users
		set user_name = $1, email = $2, updated_at = $3
		where id = $4
		`
	_, err := p.DB.ExecContext(ctx, i, u.UserName, u.Email, time.Now(), u.Id)
	if err != nil {
		log.Println("unable to update the user")
		return err
	}

	return nil
}

// AuthenticateUser compares user emain and pasword hash
func (p *postgresDBRepo) AuthenticateUser(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hash string

	m := `
		select id, password
		from users
		where email = $1
		`
	row := p.DB.QueryRowContext(ctx, m, email)
	err := row.Scan(&id, &hash)
	if err != nil {
		log.Println("unable to get user id and hash")
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hash, nil
}
