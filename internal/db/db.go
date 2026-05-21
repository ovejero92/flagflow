package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Queries {
	return &Queries{db: db}
}
