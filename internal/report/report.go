package report

import (
	"github.com/jb843051627/tephra-chron/internal/model"
	"sort"
)

func Mean(values []model.Measurement) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, v := range values {
		total += v.Value
	}
	return total / float64(len(values))
}
func SortedNotes(notes []string) []string {
	out := append([]string(nil), notes...)
	sort.Strings(out)
	return out
}
func Accepted(measurements []model.Measurement) bool {
	return len(measurements) >= 3 && Mean(measurements) > 0
}
