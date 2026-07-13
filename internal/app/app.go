package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/isaac/statescore/internal/api"
	"github.com/isaac/statescore/internal/browser"
	"github.com/isaac/statescore/internal/config"
	"github.com/isaac/statescore/internal/database"
	"github.com/isaac/statescore/internal/shutdown"
	"github.com/isaac/statescore/internal/webui"
	"github.com/isaac/statescore/web"
)

// Run starts the StateScore application.
func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	slog.Info("startup",
		"event", "startup",
		"version", config.Version,
		"dataDir", cfg.DataDir,
	)

	db, err := database.Open(cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			slog.Error("database close failed", "event", "db_close_error", "err", cerr)
		} else {
			slog.Info("database closed", "event", "db_closed")
		}
	}()

	if err := database.Migrate(db); err != nil {
		return err
	}

	ui, err := webui.New(web.Dist)
	if err != nil {
		return fmt.Errorf("load frontend assets: %w", err)
	}
	if !ui.HasAssets() {
		slog.Warn("embedded frontend missing; build frontend before go build",
			"event", "frontend_missing",
		)
	}

	mux := http.NewServeMux()
	api.NewHandler(db).Mount(mux)
	mux.Handle("/", ui)

	ln, addr, err := listenLocal(cfg.Host, cfg.Port)
	if err != nil {
		return err
	}

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := shutdown.Context()
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server started",
			"event", "server_started",
			"address", addr,
			"version", config.Version,
			"frontend", ui.HasAssets(),
		)
		fmt.Printf("\nStateScore is running.\n\nOpen: http://%s\nData: %s\n\nPress Ctrl+C to stop.\n\n", addr, cfg.DatabasePath)
		errCh <- server.Serve(ln)
	}()

	url := "http://" + addr
	if cfg.OpenBrowser {
		if err := browser.Open(url); err != nil {
			slog.Warn("browser open failed", "event", "browser_open_failed", "err", err, "url", url)
			fmt.Fprintf(os.Stderr, "Could not open the browser automatically. Open: %s\n", url)
		}
	}

	select {
	case <-ctx.Done():
		slog.Info("shutdown requested", "event", "shutdown_requested")
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	slog.Info("shutdown complete", "event", "shutdown_complete")
	return nil
}

func listenLocal(host string, startPort int) (net.Listener, string, error) {
	var lastErr error
	for port := startPort; port < startPort+50; port++ {
		addr := fmt.Sprintf("%s:%d", host, port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			lastErr = err
			continue
		}
		if port != startPort {
			slog.Info("preferred port unavailable, using fallback",
				"event", "port_fallback",
				"preferred", startPort,
				"selected", port,
			)
		}
		return ln, addr, nil
	}
	return nil, "", fmt.Errorf("no available localhost port near %d: %w", startPort, lastErr)
}
