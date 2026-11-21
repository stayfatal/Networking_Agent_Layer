package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stayfatal/Networking_Agent_Layer/internal/domain"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL, model string
	httpClient     *http.Client
}

func New(baseURL, model string) *Client {
	return &Client{baseURL: baseURL, model: model, httpClient: &http.Client{Timeout: 4 * time.Minute}}
}
func (c *Client) Generate(req domain.Request) ([]string, error) {
	body := map[string]any{"model": c.model, "stream": false, "format": map[string]any{"type": "object", "properties": map[string]any{"names": map[string]any{"type": "array", "items": map[string]string{"type": "string"}, "minItems": 5, "maxItems": 10}}, "required": []string{"names"}}, "messages": []map[string]string{{"role": "system", "content": "Ты придумываешь названия товаров на русском языке. Верни JSON с полем names: массив из 6 коротких, разных и уместных названий. Не добавляй объяснения."}, {"role": "user", "content": req.ProductDescription}}, "options": map[string]any{"temperature": 0.7, "num_predict": 500}}
	b, _ := json.Marshal(body)
	res, e := c.httpClient.Post(c.baseURL+"/api/chat", "application/json", bytes.NewReader(b))
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		v, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return nil, fmt.Errorf("ollama %d: %s", res.StatusCode, v)
	}
	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if e = json.NewDecoder(res.Body).Decode(&result); e != nil {
		return nil, e
	}
	var parsed struct {
		Names []string `json:"names"`
	}
	if e = json.Unmarshal([]byte(result.Message.Content), &parsed); e != nil {
		return nil, e
	}
	names := []string{}
	seen := map[string]bool{}
	for _, v := range parsed.Names {
		v = strings.TrimSpace(v)
		k := strings.ToLower(v)
		if v != "" && !seen[k] {
			names = append(names, v)
			seen[k] = true
		}
	}
	if len(names) < 5 {
		return nil, fmt.Errorf("model returned fewer than five unique names")
	}
	if len(names) > 10 {
		names = names[:10]
	}
	return names, nil
}
