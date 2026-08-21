package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jb843051627/tephra-chron/internal/model"
)

func (s *Store) PutChronologyReport(v model.ChronologyReport) error { return s.Save("report", v.ID, v) }
func (s *Store) GetChronologyReport(id string) (model.ChronologyReport, error) {
	var v model.ChronologyReport
	err := s.Load("report", id, &v)
	if errors.Is(err, sql.ErrNoRows) {
		return v, model.ErrMissing
	}
	return v, err
}
func (s *Store) ListChronologyReports() ([]model.ChronologyReport, error) {
	out := []model.ChronologyReport{}
	err := s.List("report", func(raw []byte) error {
		var v model.ChronologyReport
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		out = append(out, v)
		return nil
	})
	return out, err
}
