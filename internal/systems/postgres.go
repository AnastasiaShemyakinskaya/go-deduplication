package systems

import (
	"context"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewDbConn(connString string) (*Postgres, error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}
