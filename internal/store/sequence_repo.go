package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jb843051627/tephra-chron/internal/model"
)

func (s *Store) PutSequence(v model.Sequence) error { return s.Save("sequence", v.ID, v) }
func (s *Store) GetSequence(id string) (model.Sequence, error) {
	var v model.Sequence
	err := s.Load("sequence", id, &v)
	if errors.Is(err, sql.ErrNoRows) {
		return v, model.ErrMissing
	}
	return v, err
}
func (s *Store) ListSequences() ([]model.Sequence, error) {
	out := []model.Sequence{}
	err := s.List("sequence", func(raw []byte) error {
		var v model.Sequence
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		out = append(out, v)
		return nil
	})
	return out, err
}
