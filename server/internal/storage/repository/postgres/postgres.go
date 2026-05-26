package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	conn *pgx.Conn
}

func NewPostgresRepositoryFromConnection(conn *pgx.Conn) *PostgresRepository {
	repository := PostgresRepository{conn}
	return &repository
}

func NewPostgresRepository(ctx context.Context, url string) (*PostgresRepository, error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}

	return NewPostgresRepositoryFromConnection(conn), nil
}

func (repository PostgresRepository) Close(ctx context.Context) {
	repository.conn.Close(ctx)
}
