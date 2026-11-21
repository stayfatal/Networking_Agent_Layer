package transport

import (
	"encoding/json"
	"github.com/stayfatal/Networking_Agent_Layer/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSegmentBytesAndReassembly(t *testing.T) {
	var manifest domain.Manifest
	var segments []domain.Segment
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/manifest":
			json.NewDecoder(r.Body).Decode(&manifest)
		case "/api/v1/send":
			var s domain.Segment
			json.NewDecoder(r.Body).Decode(&s)
			segments = append(segments, s)
		}
		w.WriteHeader(202)
	}))
	defer srv.Close()
	client := New(srv.URL)
	req := domain.Request{Sender: "Имя", SentAt: "2026-01-01T00:00:00Z", RequestID: "test", ProductDescription: "товар"}
	data := []byte(strings.Repeat("а", 1201))
	if err := client.Send(req, "description", 0, 0, data); err != nil {
		t.Fatal(err)
	}
	if manifest.SegmentTotal != 3 || len(segments) != 3 {
		t.Fatalf("wrong segment count: %d %d", manifest.SegmentTotal, len(segments))
	}
	var got []byte
	for i, s := range segments {
		if len(s.Payload) > 1000 || s.SegmentNumber != i+1 {
			t.Fatalf("invalid segment %d", i)
		}
		got = append(got, s.Payload...)
	}
	if string(got) != string(data) {
		t.Fatal("UTF-8 bytes changed")
	}
}
