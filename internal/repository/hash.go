package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HashRepository interface {
	InsertHash(ctx context.Context, tx pgx.Tx, hash []byte) (int64, error)
	GetHash(ctx context.Context, hash []byte) (bool, error)
}

type HashRepo struct {
	db *pgxpool.Pool
}

func NewHashRepo(db *pgxpool.Pool) *HashRepo {
	return &HashRepo{db: db}
}

func (h *HashRepo) InsertHash(ctx context.Context, tx pgx.Tx, hash []byte) (int64, error) {
	var id int64
	query := `
		insert into hash (hash_string) values ($1) on conflict (hash_string) do update 
 		set hash_string = excluded.hash_string
		returning id
	`

	err := tx.QueryRow(ctx, query, hash).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (h *HashRepo) GetHash(ctx context.Context, hash []byte) (bool, error) {
	var count int64
	query := `select count(*) from hash where hash_string = $1`
	err := h.db.QueryRow(ctx, query, hash).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}
