package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func NewPostgres(cfg Config) (*Postgres, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.Username, cfg.Password, cfg.Hostname, cfg.Name)
	//slog.Debug(fmt.Sprintf("Connecting to database: %s", url))
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		conn: conn,
	}, nil
}

type Postgres struct {
	conn *pgx.Conn
}

func (p Postgres) Close() error {
	return p.conn.Close(context.Background())
}

func (p Postgres) Conn() *pgx.Conn {
	return p.conn
}
