package jobs

import (
	"context"
	"testing"
	"time"
)

func TestShutdownCancelsAndWaits(t *testing.T) {
	m := New(context.Background())
	stopped := make(chan struct{})
	m.Go(func(ctx context.Context) { <-ctx.Done(); close(stopped) })
	if m.Active() != 1 {
		t.Fatalf("active = %d, want 1", m.Active())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := m.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
	select {
	case <-stopped:
	default:
		t.Fatal("job was not canceled")
	}
}
