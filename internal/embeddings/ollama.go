package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaEmbedder struct {
	model string
}

type embedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embedResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`
}

func NewOllamaEmbedder(model string) *OllamaEmbedder {
	return &OllamaEmbedder{
		model: model,
	}
}

func (e *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	client := &http.Client{}
	payload := embedRequest{Model: e.model, Input: text}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:11434/api/embed", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var res embedResponse

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	fmt.Println("this is before returning")

	fmt.Println(res.Model)
	fmt.Println(len(res.Embeddings[0]))

	return res.Embeddings[0], nil

}
