package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type user struct{}

// Users exposes the user service
var Users user

// Create inserts a new user into the database.
func (user) Create(ctx context.Context, db *sqlx.DB, n NewUser, now time.Time) (*User, error) {
	u := User{
		ID:          uuid.New().String(),
		Name:        n.Name,
		Email:       n.Email,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO users
		(user_id, name, email, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := db.ExecContext(
		ctx, q,
		u.ID, u.Name, u.Email,
		u.DateCreated, u.DateUpdated,
	)
	if err != nil {
		return nil, errors.Wrap(err, "inserting user")
	}

	return &u, nil
}

// Retrieve gets the specified user from the database.
func (user) Retrieve(ctx context.Context, db *sqlx.DB, id string) (*User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var u User
	const q = `SELECT * FROM users WHERE user_id = $1`
	if err := db.GetContext(ctx, &u, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrapf(err, "selecting user %q", id)
	}

	return &u, nil
}
