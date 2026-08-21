package model

import "time"

type InstrumentSession struct {
	ID, SequenceID, InstrumentID, State string
	OpenedAt                            time.Time
	ClosedAt                            *time.Time
	LastHeartbeat                       time.Time
}
type Measurement struct {
	ID, SessionID      string
	Channel            string
	Value, Uncertainty float64
	CapturedAt         time.Time
	Operator           string
}

func (m Measurement) Valid() bool {
	return m.ID != "" && m.SessionID != "" && m.Channel != "" && m.Uncertainty >= 0
}
