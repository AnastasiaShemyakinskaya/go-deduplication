package entity

type FileHash struct {
	ID       int64
	HashID   int64
	FileID   int64
	Repeat   int64
	Position []int64
}
