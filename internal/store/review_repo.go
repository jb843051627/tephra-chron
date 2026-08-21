package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jb843051627/tephra-chron/internal/model"
)

func (s *Store) PutReview(v model.Review) error { return s.Save("review", v.ID, v) }
func (s *Store) GetReview(id string) (model.Review, error) {
	var v model.Review
	err := s.Load("review", id, &v)
	if errors.Is(err, sql.ErrNoRows) {
		return v, model.ErrMissing
	}
	return v, err
}
func (s *Store) ListReviews() ([]model.Review, error) {
	out := []model.Review{}
	err := s.List("review", func(raw []byte) error {
		var v model.Review
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		out = append(out, v)
		return nil
	})
	return out, err
}
