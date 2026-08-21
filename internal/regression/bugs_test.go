package regression

import (
	"context"
	"testing"

	"github.com/jb843051627/tephra-chron/internal/model"
	"github.com/jb843051627/tephra-chron/internal/service"
)

func TestBug02_CanceledQualityGate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sequence := model.Sequence{
		ID:    "sequence-02",
		State: "running",
		Steps: []model.SequenceStep{{Position: 1, Operation: "screen", State: "done"}},
	}
	_, err := (service.Gate02{}).Evaluate(ctx, sequence)
	if err == nil {
		t.Fatal("canceled laboratory request continued through a quality gate")
	}
}
