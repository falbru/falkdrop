package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	conn *pgx.Conn
}

func NewPostgresRepositoryFromConnection(conn *pgx.Conn) *PostgresRepository {
	repository := PostgresRepository{conn}
	return &repository
}

func NewPostgresRepository() (*PostgresRepository, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	return NewPostgresRepositoryFromConnection(conn), nil
}

func (repository PostgresRepository) Close() {
	repository.conn.Close(context.Background())
}
