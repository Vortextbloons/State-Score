// Package jobs manages application-owned background work.
package jobs

import (
	"context"
	"sync"
	"sync/atomic"
)

// Manager owns the lifecycle of background jobs.
type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	active atomic.Int64
}

func New(parent context.Context) *Manager {
	ctx, cancel := context.WithCancel(parent)
	return &Manager{ctx: ctx, cancel: cancel}
}

func (m *Manager) Go(job func(context.Context)) {
	m.active.Add(1)
	m.wg.Add(1)
	go func() {
		defer m.active.Add(-1)
		defer m.wg.Done()
		job(m.ctx)
	}()
}

func (m *Manager) Active() int64 { return m.active.Load() }

// Shutdown cancels jobs and waits for them, or for ctx to expire.
func (m *Manager) Shutdown(ctx context.Context) error {
	m.cancel()
	done := make(chan struct{})
	go func() { m.wg.Wait(); close(done) }()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
