package repository

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	UserName  string `json:"username"`
	Emai      string `json:"email"`
	Password  string `json:"-"` // Password is not marshed
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (user_name, email, password)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.UserName,
		user.Emai,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
