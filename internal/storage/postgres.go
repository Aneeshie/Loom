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

func (s *Store) GetLogs(ctx context.Context) ([]*pb.Log, error) {
	rows, err := s.db.Query(
		ctx,
		`
			SELECT
				id,
				service_name,
				host,
				level,
				message,
				timestamp
			FROM logs
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var logs []*pb.Log

	for rows.Next() {
		logEntry := &pb.Log{}
		err := rows.Scan(
			&logEntry.Id,
			&logEntry.ServiceName,
			&logEntry.Host,
			&logEntry.Level,
			&logEntry.Message,
			&logEntry.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, logEntry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil

}
