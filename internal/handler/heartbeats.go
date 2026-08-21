package handler

import (
	"net/http"
	"strings"
)

func (s *Server) heartbeats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("instrument"))
	if err := s.lab.Heartbeat(r.Context(), id); err != nil {
		writeJSON(w, 422, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 202, map[string]string{"instrument": id})
}
