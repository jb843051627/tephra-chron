package model

import "time"

type Sequence struct {
	ID, SpecimenID, InstrumentID, State string
	Steps                               []SequenceStep
	CreatedAt                           time.Time
	StartedAt                           *time.Time
	ClosedAt                            *time.Time
}
type SequenceStep struct {
	Position         int
	Operation, State string
	StartedAt        *time.Time
	FinishedAt       *time.Time
}

func (s Sequence) CurrentStep() *SequenceStep {
	for i := range s.Steps {
		if s.Steps[i].State != "done" {
			return &s.Steps[i]
		}
	}
	return nil
}
