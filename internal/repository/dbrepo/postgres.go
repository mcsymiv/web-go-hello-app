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
		select id, query, description, query_link, created_at, updated_at
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
		select id, query, description, query_link, created_at, updated_at
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
		(query, user_id, query_link, description, created_at, updated_at)
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
		select id, user_name, email, access_level, searches, created_at, update_at
		from users
		where id = $1
		`
	row := p.DB.QueryRowContext(ctx, i, userId)

	var u models.User

	err := row.Scan(
		&u.Id,
		&u.UserName,
		&u.Email,
		&u.AccessLevel,
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

func (p *postgresDBRepo) GetUsersCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var uc int

	q := `
		select count(*) as users_count from users	
	`
	row := p.DB.QueryRowContext(ctx, q)
	err := row.Scan(&uc)
	if err != nil {
		log.Println("unable to get users count from DB")
		return uc, err
	}

	return uc, nil

}

func (p *postgresDBRepo) GetUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []models.User

	q := `
		select u.id, u.username, u.email, u.access_level, u.created_at, u.updated_at
		from users u
	`

	rows, err := p.DB.QueryContext(ctx, q)
	if err != nil {
		log.Println("unable to get users from db", err)
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.Id,
			&u.UserName,
			&u.Email,
			&u.AccessLevel,
			&u.CreatedAt,
			&u.UpdatedAt,
		)

		if err != nil {
			log.Println("unable to scan user")
			return users, err
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		log.Println("unable to get users from db")
		return users, err
	}

	return users, nil

}

// GetSearchesByUserId returns all searches
func (p *postgresDBRepo) GetSearchesByUserId(id int) ([]models.Search, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var searches []models.Search

	q := `
		select s.id, s.query, s.query_link, s.description, s.updated_at
		from searches s
		where s.user_id = $1
	`

	rows, err := p.DB.QueryContext(ctx, q, id)
	if err != nil {
		log.Println("unable to get searches by user id from db")
		return searches, err
	}

	defer rows.Close()

	for rows.Next() {
		var s models.Search
		err := rows.Scan(
			&s.Id,
			&s.Query,
			&s.Link,
			&s.Description,
			&s.UpdatedAt,
		)

		if err != nil {
			log.Println("unable to scan searches")
			return searches, err
		}

		searches = append(searches, s)
	}

	if err = rows.Err(); err != nil {
		log.Println("unable to get searches from db")
		return searches, err
	}

	return searches, nil
}

// GetUserSearchById returns user search by user and search id
func (p *postgresDBRepo) GetUserSearchById(userId, searchId int) (models.Search, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var s models.Search

	q := `
		select s.id, s.query, s.query_link, s.description
		from searches s
		where s.user_id = $1 and s.id = $2
	`
	row := p.DB.QueryRowContext(ctx, q, userId, searchId)

	err := row.Scan(
		&s.Id,
		&s.Query,
		&s.Link,
		&s.Description,
	)

	if err != nil {
		log.Println("Unable to get search for user")
		return s, err
	}

	return s, nil
}

// UpdateUserSearch updates user search info by id
func (p *postgresDBRepo) UpdateUserSearch(s models.Search, userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	i := `
		update searches
		set query = $1, query_link = $2, description = $3, updated_at = $4
		where user_id = $5 and id = $6
		`

	_, err := p.DB.ExecContext(ctx, i, s.Query, s.Link, s.Description, time.Now(), userId, s.Id)
	if err != nil {
		log.Println("unable to update the search for user")
		return err
	}

	return nil
}

// DeleteUserSearch removes user search by id and userId
func (p *postgresDBRepo) DeleteUserSearch(searchId, userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	i := `
		delete from searches
		where user_id = $1 and id = $2
		`

	_, err := p.DB.ExecContext(ctx, i, userId, searchId)
	if err != nil {
		log.Println("unable to remove the search for user")
		return err
	}

	return nil
}

// AddUser adds new user
func (p *postgresDBRepo) AddUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	i := `
		insert into 
		users
		(username, email, password, created_at, updated_at, access_level)
		values 
		($1, $2, $3, $4, $5, $6)
	`

	_, err := p.DB.ExecContext(ctx, i, u.UserName, u.Email, u.Password, time.Now(), time.Now(), u.AccessLevel)
	if err != nil {
		log.Println("unable to insert user")
		return err
	}

	return nil
}
