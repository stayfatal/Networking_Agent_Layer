package domain

import (
	"errors"
	"strings"
	"time"
)

type Request struct {
	Sender             string `json:"sender"`
	SentAt             string `json:"sent_at"`
	RequestID          string `json:"request_id"`
	ProductDescription string `json:"product_description"`
}
type Chunk struct {
	Request
	ChunkIndex   int      `json:"chunk_index"`
	ChunkTotal   int      `json:"chunk_total"`
	IsError      bool     `json:"is_error"`
	Payload      []string `json:"payload"`
	ErrorMessage string   `json:"error_message,omitempty"`
}
type Segment struct {
	RequestID     string `json:"request_id"`
	SentAt        string `json:"sent_at"`
	Kind          string `json:"kind"`
	ChunkIndex    int    `json:"chunk_index"`
	ChunkTotal    int    `json:"chunk_total"`
	SegmentTotal  int    `json:"segment_total"`
	SegmentNumber int    `json:"segment_number"`
	Payload       []byte `json:"payload"`
}
type Manifest struct {
	Request
	Kind         string `json:"kind"`
	ChunkIndex   int    `json:"chunk_index"`
	ChunkTotal   int    `json:"chunk_total"`
	SegmentTotal int    `json:"segment_total"`
}

func (r Request) Validate() error {
	if strings.TrimSpace(r.Sender) == "" || strings.TrimSpace(r.ProductDescription) == "" || len(r.ProductDescription) > 50000 || len(r.Sender) > 100 || r.RequestID == "" {
		return errors.New("invalid request")
	}
	if _, e := time.Parse(time.RFC3339Nano, r.SentAt); e != nil {
		return errors.New("invalid sent_at")
	}
	return nil
}
