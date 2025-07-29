package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB (ctx context.Context, host string) (*DB, error) {
	connStr := "postgres://metrics:qwerty!23@"+host+"/metrics_db?sslmode=disable"

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &DB{
		conn: conn,
	}, nil
}

func (db *DB) Ping () bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
	err := db.conn.Ping(ctx)
	return err == nil
}