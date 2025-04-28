package postgres

import (
	"company/internal/repository"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func New(storagePath string) (*repository.Storage, error) {
	const op = "internal.storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &repository.Storage{Db: db}, nil
}
