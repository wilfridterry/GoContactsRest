package psql

import (
	"contact-list/internal/domain"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Users struct {
	Conn *pgx.Conn
}

func NewUsers(conn *pgx.Conn) *Users {
	return &Users{conn}
}

func (repo *Users) Create(ctx context.Context, user *domain.User) (int64, error) {
	var lastInsertId int64

	err := repo.Conn.QueryRow(
		ctx,
		"INSERT INTO users (name, email, password, registered_at) values ($1, $2, $3, $4) RETURNING id", 
		user.Name, 
		user.Email, 
		user.Password,
		user.RegisteredAt,
	).Scan(&lastInsertId)

	return lastInsertId, err
}

func (repo *Users) GetByEmailAndPassword(ctx context.Context, email string, password string) (*domain.User, error) {
	var u domain.User
	err := repo.Conn.QueryRow(ctx, "SELECT * FROM users WHERE email=$1 AND password=$2", email, password).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundUser
		}

		return nil, err
	}

	return &u, nil
}