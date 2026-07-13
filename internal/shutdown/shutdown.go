package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Context returns a context that is canceled when the process receives
// an interrupt or termination signal (Ctrl+C / SIGTERM).
func Context() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
