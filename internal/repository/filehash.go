package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type FileHashRepository interface {
	InsertFileHash(ctx context.Context, tx pgx.Tx, fileID, hashID, position int64) error
	GetPositions(ctx context.Context, file string) ([][]int64, error)
}

type FileHashRepo struct {
	db *pgxpool.Pool
}

func NewFileHashRepo(db *pgxpool.Pool) *FileHashRepo {
	return &FileHashRepo{db: db}
}

func (h *FileHashRepo) InsertFileHash(ctx context.Context, tx pgx.Tx, fileID, hashID, position int64) error {
	query := `
		insert into file_hash (hash_id, file_id, repeat, position) values ($1, $2, $3, $4) 
		on conflict (hash_id, file_id) do update 
 		set repeat = excluded.repeat + 1,
			position = array_append(file_hash.position, $5);
	`

	_, err := tx.Exec(ctx, query, hashID, fileID, 1, []int64{position}, position)
	if err != nil {
		return err
	}

	return nil
}

func (h *FileHashRepo) GetPositions(ctx context.Context, file string) ([][]int64, error) {
	positions := make([][]int64, 0)
	query := `select fh.position from file_hash fh join file f on fh.file_id = f.id where f.file_name = $1`
	rows, err := h.db.Query(ctx, query, file)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	for rows.Next() {
		position := make([]int64, 0)
		err := rows.Scan(&position)
		if err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}
	return positions, nil
}
