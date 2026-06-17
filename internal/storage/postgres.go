package storage

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/Aneeshie/loom/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

const SIMILARITY_THRESHOLD = 0.45

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

func (s *Store) InsertLog(ctx context.Context, req *pb.SendLogRequest, embedding []float32) error {

	vec := pgvector.NewVector(embedding)

	_, err := s.db.Exec(
		ctx,
		`
		INSERT INTO logs (
		service_name,
		host,
		level,
		message,
		timestamp,
		embedding
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		req.ServiceName,
		req.Host,
		req.Level,
		req.Message,
		req.Timestamp,
		vec,
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

	query := `
		SELECT
			id,
			service_name,
			host,
			level,
			message,
			timestamp
		FROM logs
	`

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

		if filter.Search != nil {
			conditions = append(conditions, fmt.Sprintf("message ILIKE $%d", len(args)+1))
			args = append(args, "%"+*filter.Search+"%")
		}

		if filter.StartTime != nil {
			conditions = append(conditions, fmt.Sprintf("timestamp >= $%d", len(args)+1))
			args = append(args, *filter.StartTime)
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	args = append(args, limit)

	query += " ORDER BY timestamp DESC"

	query += fmt.Sprintf(" LIMIT $%d", len(args))

	rows, err := s.db.Query(
		ctx,
		query,
		args...,
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

func (s *Store) SimilarLogs(ctx context.Context, embedding []float32, limit int64) ([]*pb.Log, error) {

	fmt.Println("entered similar logs")

	query := `SELECT
			id,
			service_name,
			host,
			level,
			message,
			timestamp,
			embedding <=> $1 AS distance
		FROM logs
		ORDER BY distance
		LIMIT $2`

	rows, err := s.db.Query(
		ctx,
		query,
		pgvector.NewVector(embedding),
		limit,
	)

	var distance float64

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
			&distance,
		)

		if err != nil {
			return nil, err
		}

		if distance > SIMILARITY_THRESHOLD {
			continue
		}

		fmt.Printf(
			"distance=%.4f message=%s\n",
			distance,
			logEntry.Message,
		)

		logs = append(logs, logEntry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
