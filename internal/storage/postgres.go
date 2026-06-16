package storage

import (
	"context"
	"fmt"
	"strings"

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

func (s *Store) GetLogs(ctx context.Context, limit int64, filter *pb.LogFilter) ([]*pb.Log, error) {
	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("Limit too high, should be under 100 and above 0")
	}

	query := "SELECT id, service_name, host, level, message, timestamp FROM logs"

	var args []any

	var conditions []string

	if filter != nil {
		if filter.Level != nil {
			conditions = append(conditions, fmt.Sprintf("level = $%d", len(args)+1))
			args = append(args, *filter.Level)
		}
		if filter.ServiceName != nil {
			conditions = append(conditions, fmt.Sprintf("service_name = $%d", len(args)+1))
			args = append(args, *filter.ServiceName)
		}
		if filter.Host != nil {
			conditions = append(conditions, fmt.Sprintf("host = $%d", len(args)+1))
			args = append(args, *filter.Host)
		}
	}

	if len(conditions) > 0 {
		 query += " WHERE " + strings.Join(conditions, " AND ")
	}

	args = append(args, limit)

	query += fmt.Sprintf(" LIMIT $%d", len(args))

	rows, err := s.db.Query(
		ctx,
		query,
		args,
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
