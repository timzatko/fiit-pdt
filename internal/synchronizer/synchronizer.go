package synchronizer

import (
	"context"
	"log"

	"golang.org/x/sync/semaphore"
)

// Synchronizer is used to sync inserting entities of different types
// between different files. So at one time, only to one entity is being written.
type Synchronizer struct {
	MaxWorkers int
	Ctx        context.Context

	sem *semaphore.Weighted
}

func (s *Synchronizer) Acquire() error {
	return s.sem.Acquire(s.Ctx, 1)
}

func (s *Synchronizer) Release() {
	s.sem.Release(1)
}

func (s *Synchronizer) Wait() {
	// Acquire all of the tokens to wait for any remaining workers to finish.
	if err := s.sem.Acquire(s.Ctx, int64(s.MaxWorkers)); err != nil {
		log.Printf("failed to acquire semaphore: %v", err)
	}
}

func NewSynchronizer(ctx context.Context, maxWorkers int) Synchronizer {
	var sem = semaphore.NewWeighted(int64(maxWorkers))
	return Synchronizer{Ctx: ctx, MaxWorkers: maxWorkers, sem: sem}
}
