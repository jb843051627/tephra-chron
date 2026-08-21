package report

import (
	"github.com/jb843051627/tephra-chron/internal/model"
	"sort"
)

func Timeline(items []model.Measurement) []model.Measurement {
	out := append([]model.Measurement(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i].CapturedAt.Before(out[j].CapturedAt) })
	return out
}
