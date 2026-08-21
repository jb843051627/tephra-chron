package handler

import (
	"encoding/json"
	"github.com/jb843051627/tephra-chron/internal/service"
	"net/http"
)

type Server struct {
	lab *service.Lab
	mux *http.ServeMux
}

func New(lab *service.Lab) http.Handler {
	s := &Server{lab: lab, mux: http.NewServeMux()}
	s.routes()
	return s.mux
}
func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.health)
	s.mux.HandleFunc("/api/metrics", s.metrics)
	s.mux.HandleFunc("/api/specimens", s.specimens)
	s.mux.HandleFunc("/api/sequences", s.sequences)
	s.mux.HandleFunc("/api/sessions", s.sessions)
	s.mux.HandleFunc("/api/reviews", s.reviews)
	s.mux.HandleFunc("/api/heartbeats", s.heartbeats)
	s.mux.HandleFunc("/", s.dashboard)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
