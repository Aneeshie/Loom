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

Rules:

- query: the main semantic search phrase.
- level: one of ERROR, WARN, INFO, or empty string.
- service: service name if specified.
- host: host name if specified.
- since: duration if specified (examples: 1h, 24h, 7d).
- Use empty strings for fields not mentioned.
- Do not include explanations.
- Do not include markdown.
- Do not include code fences.
- Return JSON only.

Examples:

User: show me database errors from the last hour

{
  "query": "database error",
  "level": "ERROR",
  "service": "",
  "host": "",
  "since": "1h"
}

User: show me authentication failures

{
  "query": "authentication failure",
  "level": "",
  "service": "",
  "host": "",
  "since": ""
}

User: show me warning logs from auth service

{
  "query": "",
  "level": "WARN",
  "service": "auth",
  "host": "",
  "since": ""
}

User: show me redis issues from billing service in the last day

{
  "query": "redis issue",
  "level": "",
  "service": "billing",
  "host": "",
  "since": "24h"
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
