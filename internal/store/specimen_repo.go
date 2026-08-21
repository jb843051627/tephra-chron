package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jb843051627/tephra-chron/internal/model"
)

func (s *Store) PutSpecimen(v model.Specimen) error { return s.Save("specimen", v.ID, v) }
func (s *Store) GetSpecimen(id string) (model.Specimen, error) {
	var v model.Specimen
	err := s.Load("specimen", id, &v)
	if errors.Is(err, sql.ErrNoRows) {
		return v, model.ErrMissing
	}
	return v, err
}
func (s *Store) ListSpecimens() ([]model.Specimen, error) {
	out := []model.Specimen{}
	err := s.List("specimen", func(raw []byte) error {
		var v model.Specimen
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		out = append(out, v)
		return nil
	})
	return out, err
}
