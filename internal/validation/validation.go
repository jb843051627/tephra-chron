package validation

import (
	"fmt"
	"github.com/jb843051627/tephra-chron/internal/model"
	"strings"
)

func Specimen(s model.Specimen) error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.CoreCode) == "" {
		return fmt.Errorf("specimen identity: %w", model.ErrMissing)
	}
	if s.GrainMicrons < 20 || s.GrainMicrons > 500 {
		return fmt.Errorf("grain size: %w", model.ErrRejected)
	}
	return nil
}
func Measurement(m model.Measurement) error {
	if !m.Valid() {
		return fmt.Errorf("measurement: %w", model.ErrRejected)
	}
	if m.Value != m.Value {
		return fmt.Errorf("measurement value: %w", model.ErrRejected)
	}
	return nil
}
func Transition(from, to string, allowed map[string][]string) error {
	for _, next := range allowed[from] {
		if next == to {
			return nil
		}
	}
	return fmt.Errorf("%s to %s: %w", from, to, model.ErrInvalidState)
}
