package policy

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type ApatiteCriterion struct {
	Mineral     string
	Channel     string
	Lower       float64
	Upper       float64
	DriftLimit  float64
	MinimumRuns int
}

type ApatiteObservation struct {
	Channel    string
	Value      float64
	Drift      float64
	Runs       int
	RecordedAt time.Time
}

type ApatiteAssessment struct {
	Mineral  string
	Accepted bool
	Score    float64
	Notes    []string
}

func NewApatiteCriterion() ApatiteCriterion {
	return ApatiteCriterion{
		Mineral:     "apatite",
		Channel:     "isotope-07",
		Lower:       47.0,
		Upper:       117.0,
		DriftLimit:  3.5,
		MinimumRuns: 5,
	}
}

func (c ApatiteCriterion) Validate() error {
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

func (c ApatiteCriterion) Assess(o ApatiteObservation) (ApatiteAssessment, error) {
	if err := c.Validate(); err != nil {
		return ApatiteAssessment{}, err
	}
	if o.Channel != c.Channel {
		return ApatiteAssessment{}, fmt.Errorf("channel mismatch")
	}
	notes := make([]string, 0, 4)
	score := 100.0
	if o.RecordedAt.IsZero() {
		score -= 20
		notes = append(notes, "capture time unavailable")
	}
	if math.IsNaN(o.Value) || math.IsInf(o.Value, 0) {
		return ApatiteAssessment{}, fmt.Errorf("non-finite reading")
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
	return ApatiteAssessment{Mineral: c.Mineral, Accepted: accepted, Score: score, Notes: notes}, nil
}

func (c ApatiteCriterion) Explain(o ApatiteObservation) string {
	a, err := c.Assess(o)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%s score=%.2f accepted=%t notes=%s", a.Mineral, a.Score, a.Accepted, strings.Join(a.Notes, "; "))
}
