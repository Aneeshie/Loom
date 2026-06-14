package storage

import (
	"context"

	pb "github.com/Aneeshie/loom/proto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(databaseURL string) (*Store, error) {
	conn, err := pgxpool.New(context.Background(), databaseURL)

	if err != nil {
		return nil, err
	}

	return &Store{
		db: conn,
	}, nil
}

func (s *Store) InsertLog(ctx context.Context, req *pb.SendLogRequest) error {
	_, err := s.db.Exec(
		ctx,
		`
		INSERT INTO logs (
		service_name,
		host,
		level,
		message,
		timestamp
		)
		VALUES ($1, $2, $3, $4, $5)
		`,
		req.ServiceName,
		req.Host,
		req.Level,
		req.Message,
		req.Timestamp,
	)

	return err
}

func (s *Store) CloseConnection() {
	s.db.Close()
}
