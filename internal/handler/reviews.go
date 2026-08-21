package handler

import (
	"encoding/json"
	"github.com/jb843051627/tephra-chron/internal/model"
	"net/http"
)

func (s *Server) reviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var v model.Review
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	if err := s.lab.RequestReview(r.Context(), v); err != nil {
		writeJSON(w, 422, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, v)
}
