package postgresdb

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DBConnection struct {
	pool *pgxpool.Pool
}

func NewDBConnection(ctx *context.Context, dsn string) *DBConnection {
	pool, err := pgxpool.New(*ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return &DBConnection{pool: pool}
}

func (conn *DBConnection) Close() {
	conn.pool.Close()
}

func (conn *DBConnection) GetPool() *pgxpool.Pool {
	return conn.pool
}
