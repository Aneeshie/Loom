package nlq

import "context"

type QueryIntent struct {
	Query   string `json:"query"`
	Level   string `json:"level"`
	Service string `json:"service"`
	Host    string `json:"host"`
	Since   string `json:"since"`
}

type IntentParser interface {
	Parse(
		ctx context.Context,
		query string,
	) (*QueryIntent, error)
}
