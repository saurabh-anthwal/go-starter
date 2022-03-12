package server

import (
	"context"
	"github.com/saurabh-anthwal/dummy/pkg/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var zlog = config.GetLogger()

func StartServer(cfg *config.Config) {
	zlog.Info("Initializing reauth server")
	config.InitDatabase()
	config.AutoMigrate()

	server := &http.Server{
		Addr:    cfg.GetHostPort(),
		Handler: configureRoutes(),
	}
	runCtx := signalHandlers(server)

	// Run the server
	zlog.Infof("Starting crispr server at %s", cfg.GetHostPort())
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		zlog.Fatal(err)
	}

	// Wait for server context and processors to be stopped
	<-runCtx.Done()
	zlog.Infof("crispr server gracefully exited!")
}

// os signal handlers for graceful shutdown.
// returns a context that will be cancelled on interrupts.
func signalHandlers(srvr *http.Server) context.Context {
	// Server + Processor run context
	runCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		zlog.Info("graceful shutdown initiated")

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(runCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			cancel()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				zlog.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := srvr.Shutdown(shutdownCtx)
		if err != nil {
			zlog.Fatal(err)
		}
		serverStopCtx()
	}()
	return runCtx
}

