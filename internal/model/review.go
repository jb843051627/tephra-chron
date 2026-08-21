package model

import "time"

type Review struct {
	ID, SequenceID, State, Reviewer, Note string
	RequestedAt                           time.Time
	SignedAt                              *time.Time
}
type Alert struct {
	ID, SequenceID, Kind, State, Detail string
	RaisedAt                            time.Time
	ClearedAt                           *time.Time
}

func (r Review) IsSigned() bool { return r.State == "signed" && r.SignedAt != nil }
