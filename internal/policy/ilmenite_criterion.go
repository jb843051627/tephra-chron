package policy

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type IlmeniteCriterion struct {
	Mineral     string
	Channel     string
	Lower       float64
	Upper       float64
	DriftLimit  float64
	MinimumRuns int
}

type IlmeniteObservation struct {
	Channel    string
	Value      float64
	Drift      float64
	Runs       int
	RecordedAt time.Time
}

type IlmeniteAssessment struct {
	Mineral  string
	Accepted bool
	Score    float64
	Notes    []string
}

func NewIlmeniteCriterion() IlmeniteCriterion {
	return IlmeniteCriterion{
		Mineral:     "ilmenite",
		Channel:     "isotope-09",
		Lower:       49.0,
		Upper:       119.0,
		DriftLimit:  5.5,
		MinimumRuns: 3,
	}
}

func (c IlmeniteCriterion) Validate() error {
	if strings.TrimSpace(c.Mineral) == "" {
		return fmt.Errorf("mineral is required")
	}
	if strings.TrimSpace(c.Channel) == "" {
		return fmt.Errorf("channel is required")
	}
	if c.Lower >= c.Upper {
		return fmt.Errorf("invalid calibration range")
	}
	if c.DriftLimit <= 0 {
		return fmt.Errorf("invalid drift limit")
	}
	if c.MinimumRuns < 1 {
		return fmt.Errorf("minimum runs must be positive")
	}
	return nil
}

func (c IlmeniteCriterion) Assess(o IlmeniteObservation) (IlmeniteAssessment, error) {
	if err := c.Validate(); err != nil {
		return IlmeniteAssessment{}, err
	}
	if o.Channel != c.Channel {
		return IlmeniteAssessment{}, fmt.Errorf("channel mismatch")
	}
	notes := make([]string, 0, 4)
	score := 100.0
	if o.RecordedAt.IsZero() {
		score -= 20
		notes = append(notes, "capture time unavailable")
	}
	if math.IsNaN(o.Value) || math.IsInf(o.Value, 0) {
		return IlmeniteAssessment{}, fmt.Errorf("non-finite reading")
	}
	if o.Value < c.Lower {
		score -= (c.Lower - o.Value) * 0.75
		notes = append(notes, "below calibration envelope")
	}
	if o.Value > c.Upper {
		score -= (o.Value - c.Upper) * 0.75
		notes = append(notes, "above calibration envelope")
	}
	if math.Abs(o.Drift) > c.DriftLimit {
		score -= math.Abs(o.Drift) * 3
		notes = append(notes, "instrument drift exceeds limit")
	}
	if o.Runs < c.MinimumRuns {
		score -= float64(c.MinimumRuns-o.Runs) * 8
		notes = append(notes, "insufficient repeat runs")
	}
	if score < 0 {
		score = 0
	}
	accepted := score >= 70 && len(notes) <= 1
	if len(notes) == 0 {
		notes = append(notes, "within calibration envelope")
	}
	return IlmeniteAssessment{Mineral: c.Mineral, Accepted: accepted, Score: score, Notes: notes}, nil
}

func (c IlmeniteCriterion) Explain(o IlmeniteObservation) string {
	a, err := c.Assess(o)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%s score=%.2f accepted=%t notes=%s", a.Mineral, a.Score, a.Accepted, strings.Join(a.Notes, "; "))
}
