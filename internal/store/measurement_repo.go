package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jb843051627/tephra-chron/internal/model"
)

func (s *Store) PutMeasurement(v model.Measurement) error { return s.Save("measurement", v.ID, v) }
func (s *Store) GetMeasurement(id string) (model.Measurement, error) {
	var v model.Measurement
	err := s.Load("measurement", id, &v)
	if errors.Is(err, sql.ErrNoRows) {
		return v, model.ErrMissing
	}
	return v, err
}
func (s *Store) ListMeasurements() ([]model.Measurement, error) {
	out := []model.Measurement{}
	err := s.List("measurement", func(raw []byte) error {
		var v model.Measurement
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		out = append(out, v)
		return nil
	})
	return out, err
}
