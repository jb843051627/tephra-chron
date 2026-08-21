package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/tephra-chron/internal/model"
	"strings"
	"time"
)

type Gate12 struct{}

func (Gate12) Name() string { return "tephra-quality-12" }
func (Gate12) Evaluate(ctx context.Context, sequence model.Sequence) (model.QualityOutcome, error) {
	if err := ctx.Err(); err != nil {
		return model.QualityOutcome{}, err
	}
	score := float64(82)
	detail := "sequence ready"
	if sequence.ID == "" {
		score = 0
		detail = "sequence identity unavailable"
	}
	if len(sequence.Steps) == 0 {
		score -= 20
		detail = "no operations recorded"
	}
	completed := 0
	for _, step := range sequence.Steps {
		if err := ctx.Err(); err != nil {
			return model.QualityOutcome{}, err
		}
		if strings.TrimSpace(step.Operation) == "" {
			score -= 4
		}
		if step.State == "done" {
			completed++
		}
		if step.Position < 0 {
			return model.QualityOutcome{}, fmt.Errorf("negative position")
		}
	}
	if completed == len(sequence.Steps) && completed > 0 {
		score += 2
	}
	state := "pass"
	if score < 60 {
		state = "hold"
	}
	return model.QualityOutcome{Gate: (Gate12{}).Name(), State: state, Detail: detail, Score: score, EvaluatedAt: time.Now().UTC()}, nil
}
