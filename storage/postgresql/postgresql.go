package postgresql

import (
	"app/config"
	"app/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
	user storage.UserRepoI
}

func NewConnectPostgresql(cfg *config.Config) (storage.StorageI, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%s",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))
	if err != nil {
		return nil, err
	}

	pgpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := pgpool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Store{
		db: pgpool,
		user: NewUserRepo(pgpool),
	}, nil
}

func (s *Store) CloseDb() {
	s.db.Close()
}

func (s *Store) User() storage.UserRepoI {
	if s.user == nil {
		return NewUserRepo(s.db)
	}
	return s.user
}
