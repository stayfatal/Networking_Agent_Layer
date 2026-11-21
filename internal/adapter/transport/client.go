package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stayfatal/Networking_Agent_Layer/internal/domain"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{baseURL: baseURL, httpClient: &http.Client{Timeout: 15 * time.Second}}
}
func (c *Client) Send(req domain.Request, kind string, index, total int, data []byte) error {
	n := (len(data) + 999) / 1000
	if n == 0 {
		n = 1
	}
	m := domain.Manifest{Request: req, Kind: kind, ChunkIndex: index, ChunkTotal: total, SegmentTotal: n}
	if e := c.post(c.baseURL+"/api/v1/manifest", m); e != nil {
		return e
	}
	for i := 0; i < n; i++ {
		end := (i + 1) * 1000
		if end > len(data) {
			end = len(data)
		}
		s := domain.Segment{RequestID: req.RequestID, SentAt: req.SentAt, Kind: kind, ChunkIndex: index, ChunkTotal: total, SegmentTotal: n, SegmentNumber: i + 1, Payload: data[i*1000 : end]}
		if e := c.post(c.baseURL+"/api/v1/send", s); e != nil {
			return e
		}
	}
	return nil
}
func (c *Client) post(url string, v any) error {
	b, e := json.Marshal(v)
	if e != nil {
		return e
	}
	res, e := c.httpClient.Post(url, "application/json", bytes.NewReader(b))
	if e != nil {
		return e
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return fmt.Errorf("%s: HTTP %d: %s", url, res.StatusCode, b)
	}
	return nil
}
func (c *Client) Failure(req domain.Request, err error) error {
	return c.post(c.baseURL+"/api/v1/failure", domain.Chunk{Request: req, IsError: true, ErrorMessage: "Ошибка генерации или передачи: " + err.Error()})
}
