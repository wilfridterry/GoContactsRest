package psql

import (
	"contact-list/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type Users struct {
	Conn *pgx.Conn
}

func NewUsers(conn *pgx.Conn) *Users {
	return &Users{conn}
}

func (repo *Users) Create(ctx context.Context, user domain.User) error {
	_, err := repo.Conn.Exec(
		ctx,
		"INSERT INTO users (name, email, password, registered_at) values ($1, $2, $3, $4)", 
		user.Name, 
		user.Email, 
		user.Password,
		user.RegisteredAt,
	)

	return err
}