package model

import "time"

type Specimen struct {
	ID, CoreCode, LayerCode, State string
	CollectedAt                    time.Time
	GrainMicrons                   int
	Notes                          string
}
type GrainSplit struct {
	ID, SpecimenID string
	Mesh           int
	Fraction       float64
	PreparedAt     time.Time
}

func (s Specimen) ReadyForSequence() bool {
	return s.ID != "" && s.CoreCode != "" && s.GrainMicrons > 0 && s.State == "prepared"
}
