package nlq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type OllamaParser struct {
	model string
}

func NewOllamaParser(model string) *OllamaParser {
	return &OllamaParser{
		model: model,
	}
}

func (p *OllamaParser) Parse(ctx context.Context, query string) (*QueryIntent, error) {

	basePrompt := `You are a query parser for Loom.

Your job is to convert a user's natural language request into JSON.

Return ONLY valid JSON.

Schema:

{
  "query": "",
  "level": "",
  "service": "",
  "host": "",
  "since": ""
}

Field meanings:

- query: the main subject being searched for.
- level: ERROR, WARN, INFO, or empty string.
- service: service name if specified.
- host: host name if specified.
- since: duration such as 1h, 24h, 7d.

Rules:

1. Return ONLY JSON.
2. Do not include explanations.
3. Do not include markdown.
4. Do not include code fences.
5. Use empty strings for unspecified fields.

Log Level Rules:

- "error", "errors", "failure", "failures", "failed", "exception", "exceptions"
  => level = "ERROR"

- "warning", "warnings", "warn"
  => level = "WARN"

- "info", "informational"
  => level = "INFO"

Query Rules:

- Remove log-level words from the query.
- Keep only the subject being searched.
- Prefer concise search phrases.

Examples:

User:
show me authentication failures

{
  "query": "authentication",
  "level": "ERROR",
  "service": "",
  "host": "",
  "since": ""
}

User:
show me database errors from the last hour

{
  "query": "database",
  "level": "ERROR",
  "service": "",
  "host": "",
  "since": "1h"
}

User:
show me redis issues from billing service

{
  "query": "redis",
  "level": "",
  "service": "billing",
  "host": "",
  "since": ""
}

User:
show me warning logs from auth service

{
  "query": "",
  "level": "WARN",
  "service": "auth",
  "host": "",
  "since": ""
}

User:
show me logs from host prod-1

{
  "query": "",
  "level": "",
  "service": "",
  "host": "prod-1",
  "since": ""
}

Now convert the following request:

{{QUERY}}`

	prompt := strings.ReplaceAll(
		basePrompt,
		"{{QUERY}}",
		query,
	)

	resp, err := p.callOllamaAPI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	fmt.Println("RAW OLLAMA RESPONSE:")
	fmt.Println(resp)

	var intent QueryIntent

	if err := json.Unmarshal(
		[]byte(resp),
		&intent,
	); err != nil {
		return nil, err
	}

	return &intent, nil
}

func (p *OllamaParser) callOllamaAPI(
	ctx context.Context,
	prompt string,
) (string, error) {

	client := &http.Client{}

	payload := map[string]any{
		"model":  p.model,
		"prompt": prompt,
		"stream": false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"http://localhost:11434/api/generate",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Response, nil
}

func callOllamaAPI(ctx context.Context, prompt string) (string, error) {
	client := &http.Client{}
	payload := map[string]string{
		"model":  "llama3.2:3b",
		"prompt": prompt,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:11434/api/generate", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Output string `json:"output"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Output, nil
}
