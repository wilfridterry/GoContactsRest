package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ConnectionConfig struct {
	Host     string
	Port     uint16
	Database string
	Username string
	Password string
	SSLMode  bool
}

func NewConnection(ctx context.Context, cf *ConnectionConfig) (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cf.Username, cf.Password, cf.Host, cf.Port, cf.Database)
	conn, err := pgx.Connect(ctx, connString)

	if err != nil {
		return nil, err
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
