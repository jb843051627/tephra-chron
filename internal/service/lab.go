package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/tephra-chron/internal/ingest"
	"github.com/jb843051627/tephra-chron/internal/metrics"
	"github.com/jb843051627/tephra-chron/internal/model"
	"github.com/jb843051627/tephra-chron/internal/store"
	"github.com/jb843051627/tephra-chron/internal/validation"
	"sync"
	"time"
)

type Lab struct {
	store     *store.Store
	queue     *ingest.Queue
	metrics   *metrics.Registry
	mu        sync.RWMutex
	heartbeat map[string]time.Time
	gates     []QualityGate
}
type QualityGate interface {
	Name() string
	Evaluate(context.Context, model.Sequence) (model.QualityOutcome, error)
}

func NewLab(s *store.Store) *Lab {
	app := &Lab{store: s, queue: ingest.New(32), metrics: metrics.New(), heartbeat: map[string]time.Time{}}
	app.gates = allGates()
	return app
}
func (l *Lab) Close() { l.queue.Close() }
func (l *Lab) RegisterSpecimen(ctx context.Context, v model.Specimen) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validation.Specimen(v); err != nil {
		return err
	}
	v.State = "prepared"
	if err := l.store.PutSpecimen(v); err != nil {
		return fmt.Errorf("save specimen: %w", err)
	}
	l.metrics.Add("specimens", 1)
	return l.store.Event(v.ID, "registered")
}
func (l *Lab) ScheduleSequence(ctx context.Context, v model.Sequence) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	sp, err := l.store.GetSpecimen(v.SpecimenID)
	if err != nil {
		return fmt.Errorf("sequence specimen: %w", err)
	}
	if !sp.ReadyForSequence() {
		return model.ErrInvalidState
	}
	v.State = "queued"
	v.CreatedAt = time.Now().UTC()
	if err := l.store.PutSequence(v); err != nil {
		return err
	}
	return l.store.Event(v.ID, "scheduled")
}
func (l *Lab) OpenSession(ctx context.Context, v model.InstrumentSession) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	seq, err := l.store.GetSequence(v.SequenceID)
	if err != nil {
		return err
	}
	if seq.State != "queued" {
		return model.ErrInvalidState
	}
	now := time.Now().UTC()
	v.State = "open"
	v.OpenedAt = now
	v.LastHeartbeat = now
	seq.State = "running"
	seq.StartedAt = &now
	if err := transactionMarker(l.store, v.SequenceID); err != nil {
		return err
	}
	if err := l.store.PutSession(v); err != nil {
		return err
	}
	return l.store.PutSequence(seq)
}
func (l *Lab) RecordMeasurement(ctx context.Context, v model.Measurement) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validation.Measurement(v); err != nil {
		return err
	}
	session, err := l.store.GetInstrumentSession(v.SessionID)
	if err != nil {
		return err
	}
	if session.State != "open" {
		return model.ErrInvalidState
	}
	v.CapturedAt = time.Now().UTC()
	if err := l.store.PutMeasurement(v); err != nil {
		return err
	}
	l.metrics.Add("measurements", 1)
	return l.store.Event(v.SessionID, "measurement")
}
func (l *Lab) CloseSession(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	v, err := l.store.GetInstrumentSession(id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	v.State = "closed"
	v.ClosedAt = &now
	if err := l.store.PutInstrumentSession(v); err != nil {
		return err
	}
	seq, err := l.store.GetSequence(v.SequenceID)
	if err != nil {
		return err
	}
	seq.State = "review"
	return l.store.PutSequence(seq)
}
func (l *Lab) RequestReview(ctx context.Context, v model.Review) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	seq, err := l.store.GetSequence(v.SequenceID)
	if err != nil {
		return err
	}
	if seq.State != "review" {
		return model.ErrInvalidState
	}
	v.State = "requested"
	v.RequestedAt = time.Now().UTC()
	return l.store.PutReview(v)
}
func (l *Lab) SignReview(ctx context.Context, id, reviewer string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	v, err := l.store.GetReview(id)
	if err != nil {
		return err
	}
	if v.State != "requested" {
		return model.ErrInvalidState
	}
	now := time.Now().UTC()
	v.State = "signed"
	v.Reviewer = reviewer
	v.SignedAt = &now
	if err := l.store.PutReview(v); err != nil {
		return err
	}
	seq, err := l.store.GetSequence(v.SequenceID)
	if err != nil {
		return err
	}
	seq.State = "accepted"
	return l.store.PutSequence(seq)
}
func (l *Lab) Heartbeat(ctx context.Context, instrument string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	l.mu.Lock()
	l.heartbeat[instrument] = time.Now().UTC()
	l.mu.Unlock()
	l.metrics.Add("heartbeats", 1)
	return nil
}
func (l *Lab) Evaluate(ctx context.Context, sequenceID string) ([]model.QualityOutcome, error) {
	seq, err := l.store.GetSequence(sequenceID)
	if err != nil {
		return nil, err
	}
	out := make([]model.QualityOutcome, 0, len(l.gates))
	for _, gate := range l.gates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		v, err := gate.Evaluate(ctx, seq)
		if err != nil {
			return nil, fmt.Errorf("gate %s: %w", gate.Name(), err)
		}
		out = append(out, v)
	}
	return out, nil
}
func (l *Lab) Metrics() map[string]int64 { return l.metrics.Snapshot() }
func (l *Lab) GetSequence(ctx context.Context, id string) (model.Sequence, error) {
	if err := ctx.Err(); err != nil {
		return model.Sequence{}, err
	}
	return l.store.GetSequence(id)
}
