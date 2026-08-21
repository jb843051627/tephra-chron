package service

import "github.com/jb843051627/tephra-chron/internal/store"

func transactionMarker(s *store.Store, subject string) error {
	return s.Event(subject, "state transition")
}
