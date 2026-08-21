package model

import "time"

type ChronologyReport struct {
	ID, SequenceID  string
	GeneratedAt     time.Time
	Accepted        bool
	MeanAge, Spread float64
	Notes           []string
}
type QualityOutcome struct {
	Gate, State, Detail string
	Score               float64
	EvaluatedAt         time.Time
}
