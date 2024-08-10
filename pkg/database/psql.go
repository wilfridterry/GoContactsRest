package database

import (
	"context"

	"github.com/jackc/pgx"
)

type ConnectionInfo struct {
	Host     string
	Port     uint16
	DBName   string
	Username string
	Password string
	SSLMode  bool
}

func NewConnection(info *ConnectionInfo) (*pgx.Conn, error) {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     info.Host,
		Port:     info.Port,
		Database: info.DBName,
		User:     info.Username,
		Password: info.Password,
	})

	if err != nil {
		return nil, err
	}

	if err = conn.Ping(context.TODO()); err != nil {
		return nil, err
	}

	return conn, nil
}
