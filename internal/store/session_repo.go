package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jb843051627/tephra-chron/internal/model"
)

func (s *Store) PutInstrumentSession(v model.InstrumentSession) error {
	return s.Save("session", v.ID, v)
}
func (s *Store) GetInstrumentSession(id string) (model.InstrumentSession, error) {
	var v model.InstrumentSession
	err := s.Load("session", id, &v)
	if errors.Is(err, sql.ErrNoRows) {
		return v, model.ErrMissing
	}
	return v, err
}
func (s *Store) ListInstrumentSessions() ([]model.InstrumentSession, error) {
	out := []model.InstrumentSession{}
	err := s.List("session", func(raw []byte) error {
		var v model.InstrumentSession
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		out = append(out, v)
		return nil
	})
	return out, err
}
