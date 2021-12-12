package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"go-deduplication/internal/entity"
)

type FileRepository interface {
	InsertFile(ctx context.Context, hashFunc entity.HashFunction, bytes int, file string) (int64, error)
}

type FileRepo struct {
	db *pgxpool.Pool
}

func NewFileRepo(db *pgxpool.Pool) *HashRepo {
	return &HashRepo{db: db}
}

func (h *HashRepo) InsertFile(ctx context.Context, hashFunc entity.HashFunction, bytes int, file string) (int64, error) {
	var id int64
	query := `
		insert into file (hash_function, byte_size, file_name) values ($1, $2, $3) 
		on conflict (file_name) do update set file_name=EXCLUDED.file_name 
		returning id
	`

	err := h.db.QueryRow(ctx, query, hashFunc, bytes, file).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
