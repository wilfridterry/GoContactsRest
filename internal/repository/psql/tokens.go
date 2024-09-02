package psql

import (
	"contact-list/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type Tokens struct {
	Conn *pgx.Conn
}

func NewTokens(conn *pgx.Conn) *Tokens {
	return &Tokens{conn}
}

func (r *Tokens) Create(ctx context.Context, session *domain.RefreshSession) error {
	_, err := r.Conn.Exec(
		ctx,
		"INSERT INTO refresh_tokens (user_id, token expires_at) values ($1, $2, $3)",
		session.UserId,
		session.Token,
		session.ExpiresAt,
	)

	return err
}

func (r *Tokens) GetByToken(ctx context.Context, token string) (*domain.RefreshSession, error) {
	s := domain.RefreshSession{}
	row := r.Conn.QueryRow(ctx, "SELECT * from refresh_tokens WHERE token = $1", token)

	if err := row.Scan(&s.ID, &s.Token, &s.ExpiresAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return &domain.RefreshSession{}, err
	}

	_, err := r.Conn.Exec(ctx, "DELETE FROM contacts WHERE token = $1", token)


	return &s, err
}