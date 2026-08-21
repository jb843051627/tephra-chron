package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jb843051627/tephra-chron/internal/model"
)

func (s *Store) PutAlert(v model.Alert) error { return s.Save("alert", v.ID, v) }
func (s *Store) GetAlert(id string) (model.Alert, error) {
	var v model.Alert
	err := s.Load("alert", id, &v)
	if errors.Is(err, sql.ErrNoRows) {
		return v, model.ErrMissing
	}
	return v, err
}
func (s *Store) ListAlerts() ([]model.Alert, error) {
	out := []model.Alert{}
	err := s.List("alert", func(raw []byte) error {
		var v model.Alert
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		out = append(out, v)
		return nil
	})
	return out, err
}
