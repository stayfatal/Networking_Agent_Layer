package httpapi

import (
	"encoding/json"
	"github.com/stayfatal/Networking_Agent_Layer/internal/domain"
	"io"
	"net/http"
)

type Processor interface{ Process(domain.Request) }

func New(processor Processor) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/process", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			respond(w, 405, map[string]string{"code": "error", "message": "method not allowed"})
			return
		}
		defer r.Body.Close()
		var req domain.Request
		if json.NewDecoder(io.LimitReader(r.Body, 256*1024)).Decode(&req) != nil || req.Validate() != nil {
			respond(w, 400, map[string]string{"code": "error", "message": "invalid request"})
			return
		}
		go processor.Process(req)
		respond(w, 202, map[string]string{"request_id": req.RequestID, "status": "accepted"})
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	return mux
}
func respond(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
