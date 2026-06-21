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

	condition, args := buildConditions(filter, 1)

	if len(condition) > 0 {
		query += " WHERE " + strings.Join(condition, " AND ")
	}

	args = append(args, limit)

	query += fmt.Sprintf(
		" ORDER BY timestamp DESC LIMIT $%d",
		len(args),
	)

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

func (s *Store) SimilarLogs(
	ctx context.Context,
	embedding []float32,
	limit int64,
	filter *pb.LogFilter,
) ([]*pb.Log, error) {

	query := `
		SELECT
			id,
			service_name,
			host,
			level,
			message,
			timestamp,
			embedding <=> $1 AS distance
		FROM logs
	`

	conditions, filterArgs := buildConditions(filter, 2)

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY distance"

	args := []any{
		pgvector.NewVector(embedding),
	}

	args = append(args, filterArgs...)
	args = append(args, limit)

	query += fmt.Sprintf(
		" LIMIT $%d",
		len(args),
	)

	rows, err := s.db.Query(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var distance float64
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

func buildConditions(
	filter *pb.LogFilter,
	startArg int,
) ([]string, []any) {

	var conditions []string
	var args []any

	if filter == nil {
		return conditions, args
	}

	argNum := startArg

	if filter.Level != nil {
		conditions = append(
			conditions,
			fmt.Sprintf("level = $%d", argNum),
		)
		args = append(args, *filter.Level)
		argNum++
	}

	if filter.ServiceName != nil {
		conditions = append(
			conditions,
			fmt.Sprintf("service_name = $%d", argNum),
		)
		args = append(args, *filter.ServiceName)
		argNum++
	}

	if filter.Host != nil {
		conditions = append(
			conditions,
			fmt.Sprintf("host = $%d", argNum),
		)
		args = append(args, *filter.Host)
		argNum++
	}

	if filter.Search != nil {
		conditions = append(
			conditions,
			fmt.Sprintf("message ILIKE $%d", argNum),
		)
		args = append(args, "%"+*filter.Search+"%")
		argNum++
	}

	if filter.StartTime != nil {
		conditions = append(
			conditions,
			fmt.Sprintf("timestamp >= $%d", argNum),
		)
		args = append(args, *filter.StartTime)
		argNum++
	}

	return conditions, args
}
